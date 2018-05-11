// +build linux

/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cm

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// getResourceList returns a ResourceList with the
// specified cpu and memory resource values
func getResourceList(cpu, memory string) v1.ResourceList {
	res := v1.ResourceList{}
	if cpu != "" {
		res[v1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[v1.ResourceMemory] = resource.MustParse(memory)
	}
	return res
}

// getResourceRequirements returns a ResourceRequirements object
func getResourceRequirements(requests, limits v1.ResourceList) v1.ResourceRequirements {
	res := v1.ResourceRequirements{}
	res.Requests = requests
	res.Limits = limits
	return res
}

func TestResourceConfigForPod(t *testing.T) {
	minShares := uint64(MinShares)
	burstableShares := MilliCPUToShares(100)
	memoryQuantity := resource.MustParse("200Mi")
	burstableMemory := memoryQuantity.Value()
	burstablePartialShares := MilliCPUToShares(200)
	burstablePeriod := DefaultQuotaPeriod
	burstableQuota := MilliCPUToQuota(200, burstablePeriod)
	guaranteedShares := MilliCPUToShares(100)
	guaranteedPeriod := uint64(10000)
	guaranteedQuota := MilliCPUToQuota(100, guaranteedPeriod)
	memoryQuantity = resource.MustParse("100Mi")
	cpuNoLimit := int64(-1)
	guaranteedMemory := memoryQuantity.Value()
	testCases := map[string]struct {
		pod              *v1.Pod
		expected         *ResourceConfig
		enforceCPULimits bool
	}{
		"besteffort": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{},
								Limits:   v1.ResourceList{},
							},
						},
					},
				},
			},
			enforceCPULimits: true,
			expected:         &ResourceConfig{CpuShares: &minShares},
		},
		"burstable-no-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{},
							},
						},
					},
				},
			},
			enforceCPULimits: true,
			expected:         &ResourceConfig{CpuShares: &burstableShares},
		},
		"burstable-with-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
			enforceCPULimits: true,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &burstableQuota, CpuPeriod: &burstablePeriod, Memory: &burstableMemory},
		},
		"burstable-with-limits-no-cpu-enforcement": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
			enforceCPULimits: false,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &cpuNoLimit, CpuPeriod: &burstablePeriod, Memory: &burstableMemory},
		},
		"burstable-partial-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{},
							},
						},
					},
				},
			},
			enforceCPULimits: true,
			expected:         &ResourceConfig{CpuShares: &burstablePartialShares},
		},
		"guaranteed": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:         resource.MustParse("100m"),
									v1.ResourceMemory:      resource.MustParse("100Mi"),
									"monzo.com/cpu-period": resource.MustParse("10000"),
								},
							},
						},
					},
				},
			},
			enforceCPULimits: true,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &guaranteedQuota, CpuPeriod: &guaranteedPeriod, Memory: &guaranteedMemory},
		},
		"guaranteed-no-cpu-enforcement": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:         resource.MustParse("100m"),
									v1.ResourceMemory:      resource.MustParse("100Mi"),
									"monzo.com/cpu-period": resource.MustParse("10000"),
								},
							},
						},
					},
				},
			},
			enforceCPULimits: false,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &guaranteedPeriod, Memory: &guaranteedMemory},
		},
	}
	for testName, testCase := range testCases {
		actual := ResourceConfigForPod(testCase.pod, testCase.enforceCPULimits)
		if !reflect.DeepEqual(actual.CpuPeriod, testCase.expected.CpuPeriod) {
			t.Errorf("unexpected result, test: %v, cpu period (%d) not as expected (%d)", testName, *(actual.CpuPeriod), *(testCase.expected.CpuPeriod))
		}
		if !reflect.DeepEqual(actual.CpuQuota, testCase.expected.CpuQuota) {
			t.Errorf("unexpected result, test: %v, cpu quota (%d) not as expected (%d)", testName, *(actual.CpuQuota), *(testCase.expected.CpuQuota))
		}
		if !reflect.DeepEqual(actual.CpuShares, testCase.expected.CpuShares) {
			t.Errorf("unexpected result, test: %v, cpu shares (%d) not as expected (%d)", testName, *(actual.CpuShares), *(testCase.expected.CpuShares))
		}
		if !reflect.DeepEqual(actual.Memory, testCase.expected.Memory) {
			t.Errorf("unexpected result, test: %v, memory (%d) not as expected (%d)", testName, *(actual.Memory), *(testCase.expected.Memory))
		}
	}
}

func TestMilliCPUToQuota(t *testing.T) {
	testCases := []struct {
		cpu    int64
		period uint64
		quota  int64
	}{
		{
			cpu:    int64(0),
			period: uint64(100000),
			quota:  int64(0),
		},
		{
			cpu:    int64(5),
			period: uint64(100000),
			quota:  int64(1000),
		},
		{
			cpu:    int64(9),
			period: uint64(100000),
			quota:  int64(1000),
		},
		{
			cpu:    int64(10),
			period: uint64(100000),
			quota:  int64(1000),
		},
		{
			cpu:    int64(200),
			period: uint64(100000),
			quota:  int64(20000),
		},
		{
			cpu:    int64(500),
			period: uint64(100000),
			quota:  int64(50000),
		},
		{
			cpu:    int64(1000),
			period: uint64(100000),
			quota:  int64(100000),
		},
		{
			cpu:    int64(1500),
			period: uint64(100000),
			quota:  int64(150000),
		},
		{
			cpu:    int64(1500),
			period: uint64(10000),
			quota:  int64(15000),
		},
		{
			cpu:    int64(250),
			period: uint64(5000),
			quota:  int64(1250),
		},
	}
	for _, testCase := range testCases {
		quota := MilliCPUToQuota(testCase.cpu, testCase.period)
		if quota != testCase.quota {
			t.Errorf("Input (cpu=%d, period=%d), expected quota=%d but got quota=%d", testCase.cpu, testCase.period, testCase.quota, quota)
		}
	}
}
