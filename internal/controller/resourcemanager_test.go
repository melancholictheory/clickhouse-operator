package controller

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("addedVolumeClaimTemplates", func() {
	sts := func(names ...string) *appsv1.StatefulSet {
		templates := make([]corev1.PersistentVolumeClaim, 0, len(names))
		for _, name := range names {
			templates = append(templates, corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: name},
			})
		}

		return &appsv1.StatefulSet{Spec: appsv1.StatefulSetSpec{VolumeClaimTemplates: templates}}
	}

	It("reports nothing when the disks are unchanged", func() {
		Expect(addedVolumeClaimTemplates(sts("data", "disk1"), sts("data", "disk1"))).To(BeEmpty())
	})

	It("reports the disks that the desired StatefulSet adds", func() {
		Expect(addedVolumeClaimTemplates(sts("data"), sts("data", "disk1", "disk2"))).
			To(ConsistOf("disk1", "disk2"))
	})

	It("ignores disks that only the existing StatefulSet has", func() {
		Expect(addedVolumeClaimTemplates(sts("data", "disk1"), sts("data"))).To(BeEmpty())
	})

	It("reports the first disk added to a StatefulSet that had only the data volume", func() {
		Expect(addedVolumeClaimTemplates(sts("data"), sts("data", "disk1"))).To(ConsistOf("disk1"))
	})

	It("reports nothing when either StatefulSet is missing", func() {
		Expect(addedVolumeClaimTemplates(nil, sts("data", "disk1"))).To(BeEmpty())
		Expect(addedVolumeClaimTemplates(sts("data"), nil)).To(BeEmpty())
	})
})
