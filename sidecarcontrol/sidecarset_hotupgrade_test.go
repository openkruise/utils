/*
Copyright 2020 The Kruise Authors.

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

package sidecarcontrol

import (
	"testing"

	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
)

func TestInjectHotUpgradeSidecar(t *testing.T) {
	sidecarSetIn := sidecarSet1.DeepCopy()
	sidecarSetIn.Annotations[SidecarSetHashWithoutImageAnnotation] = "without-c4k2dbb95d"
	sidecarSetIn.Spec.Containers[0].UpgradeStrategy.UpgradeType = appsv1alpha1.SidecarContainerHotUpgrade
	sidecarSetIn.Spec.Containers[0].UpgradeStrategy.HotUpgradeEmptyImage = "busy:hotupgrade-empty"
	testInjectHotUpgradeSidecar(t, sidecarSetIn)
}

func testInjectHotUpgradeSidecar(t *testing.T, sidecarSetIn *appsv1alpha1.SidecarSet) {
	podIn := pod1.DeepCopy()
	podOut := podIn.DeepCopy()
	_, err := SidecarSetMutatingPod(podOut, nil, []*appsv1alpha1.SidecarSet{sidecarSetIn}, NewCommonControl(nil, ""))
	if err != nil {
		t.Fatalf("inject sidecar into pod failed, err: %v", err)
	}
	if len(podOut.Spec.Containers) != 4 {
		t.Fatalf("expect 4 containers but got %v", len(podOut.Spec.Containers))
	}
	if podOut.Spec.Containers[0].Image != sidecarSetIn.Spec.Containers[0].Image {
		t.Fatalf("expect image %v but got %v", sidecarSetIn.Spec.Containers[0].Image, podOut.Spec.Containers[0].Image)
	}
	if podOut.Spec.Containers[1].Image != sidecarSetIn.Spec.Containers[0].UpgradeStrategy.HotUpgradeEmptyImage {
		t.Fatalf("expect image busy:hotupgrade-empty but got %v", podOut.Spec.Containers[1].Image)
	}
	if GetPodSidecarSetRevision("sidecarset1", podOut) != GetSidecarSetRevision(sidecarSetIn) {
		t.Fatalf("pod sidecarset revision(%s) error", GetPodSidecarSetRevision("sidecarset1", podOut))
	}
	if GetPodSidecarSetWithoutImageRevision("sidecarset1", podOut) != GetSidecarSetWithoutImageRevision(sidecarSetIn) {
		t.Fatalf("pod sidecarset without image revision(%s) error", GetPodSidecarSetWithoutImageRevision("sidecarset1", podOut))
	}
	if podOut.Annotations[SidecarSetListAnnotation] != "sidecarset1" {
		t.Fatalf("pod annotations[%s]=%s error", SidecarSetListAnnotation, podOut.Annotations[SidecarSetListAnnotation])
	}
	if GetPodHotUpgradeInfoInAnnotations(podOut)["dns-f"] != "dns-f-1" {
		t.Fatalf("pod annotations[%s]=%s error", SidecarSetWorkingHotUpgradeContainer, podOut.Annotations[SidecarSetWorkingHotUpgradeContainer])
	}
	if podOut.Annotations[GetPodSidecarSetVersionAnnotation("dns-f-1")] != "1" {
		t.Fatalf("pod annotations dns-f-1 version=%s", podOut.Annotations[GetPodSidecarSetVersionAnnotation("dns-f-1")])
	}
	if podOut.Annotations[GetPodSidecarSetVersionAnnotation("dns-f-2")] != "0" {
		t.Fatalf("pod annotations dns-f-2 version=%s", podOut.Annotations[GetPodSidecarSetVersionAnnotation("dns-f-2")])
	}
}
