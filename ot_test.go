package ot_test

import (
	"github.com/alexdavid/ot"
	. "github.com/alexdavid/ot/operation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go OT", func() {

	Describe("Equals", func() {
		It("returns true for two empty text operations", func() {
			Expect(ot.Equals(
				new(ot.Transform),
				new(ot.Transform),
			)).To(BeTrue())
		})

		It("returns true if two operations have the same operations", func() {
			Expect(ot.Equals(
				new(ot.Transform).Insert("Foo"),
				new(ot.Transform).Insert("Foo"),
			)).To(BeTrue())
			Expect(ot.Equals(
				new(ot.Transform).Delete("Bar"),
				new(ot.Transform).Delete("Bar"),
			)).To(BeTrue())
			Expect(ot.Equals(
				new(ot.Transform).Retain(5),
				new(ot.Transform).Retain(5),
			)).To(BeTrue())
			Expect(ot.Equals(
				new(ot.Transform).Insert("Foo").Retain(5).Delete("Bar"),
				new(ot.Transform).Insert("Foo").Retain(5).Delete("Bar"),
			)).To(BeTrue())
		})

		It("returns false if two operations have different operation content", func() {
			Expect(ot.Equals(
				new(ot.Transform).Insert("Foo"),
				new(ot.Transform).Insert("Bar"),
			)).To(BeFalse())
			Expect(ot.Equals(
				new(ot.Transform).Delete("Fizz"),
				new(ot.Transform).Delete("Buzz"),
			)).To(BeFalse())
			Expect(ot.Equals(
				new(ot.Transform).Retain(5),
				new(ot.Transform).Retain(9),
			)).To(BeFalse())
			Expect(ot.Equals(
				new(ot.Transform).Retain(5).Insert("Foo").Delete("Bar"),
				new(ot.Transform).Insert("Foo").Retain(5).Delete("Bar"),
			)).To(BeFalse())
		})

		It("returns false if two operations have different operation types", func() {
			Expect(ot.Equals(
				new(ot.Transform).Insert("Foo"),
				new(ot.Transform).Delete("Foo"),
			)).To(BeFalse())
		})

		It("returns false if two operations have different base lengths", func() {
			Expect(ot.Equals(
				&ot.Transform{BaseLength: 1},
				&ot.Transform{BaseLength: 2},
			)).To(BeFalse())
		})

		It("returns false if two operations have different target lengths", func() {
			Expect(ot.Equals(
				&ot.Transform{TargetLength: 1},
				&ot.Transform{TargetLength: 2},
			)).To(BeFalse())
		})

	})

	Describe("Retain", func() {
		It("Does nothing when given a number less than 1", func() {
			Expect(new(ot.Transform).Retain(0).Operations()).To(Equal([]Operation{}))
			Expect(new(ot.Transform).Retain(-3).Operations()).To(Equal([]Operation{}))
		})

		It("Adds to the base length", func() {
			Expect(new(ot.Transform).Retain(5).BaseLength).To(Equal(5))
			Expect(new(ot.Transform).Retain(3).Retain(4).BaseLength).To(Equal(7))
		})

		It("Adds to the target length", func() {
			Expect(new(ot.Transform).Retain(5).TargetLength).To(Equal(5))
			Expect(new(ot.Transform).Retain(3).Retain(4).TargetLength).To(Equal(7))
		})

		It("Saves the operation", func() {
			Expect(
				new(ot.Transform).Retain(5).Operations(),
			).To(Equal([]Operation{RetainOperation{Length: 5}}))
		})

		It("Combines with previous Retain operations", func() {
			Expect(
				new(ot.Transform).Retain(5).Retain(4).Operations(),
			).To(Equal([]Operation{RetainOperation{Length: 9}}))
		})
	})

	Describe("Insert", func() {
		It("Does nothing when given an empty string", func() {
			Expect(new(ot.Transform).Insert("").Operations()).To(Equal([]Operation{}))
		})

		It("Does not change the base length", func() {
			Expect(new(ot.Transform).Insert("Foo").BaseLength).To(Equal(0))
			Expect(new(ot.Transform).Retain(3).Insert("Fizz").BaseLength).To(Equal(3))
		})

		It("Adds to the target length", func() {
			Expect(new(ot.Transform).Insert("Foo").TargetLength).To(Equal(3))
			Expect(new(ot.Transform).Insert("Foo").Insert("Fizz").TargetLength).To(Equal(7))
		})

		It("Saves the operation", func() {
			Expect(new(ot.Transform).Insert("Fizz").Operations()).To(Equal([]Operation{InsertOperation{Content: "Fizz"}}))
		})

		It("Combines with previous Insert operations", func() {
			Expect(new(ot.Transform).Insert("Foo").Insert("Bar").Operations()).To(Equal([]Operation{InsertOperation{Content: "FooBar"}}))
		})

		It("Always comes before delete operations", func() {
			Expect(new(ot.Transform).Retain(4).Delete("Foo").Insert("Bar").Operations()).To(Equal(
				[]Operation{
					RetainOperation{Length: 4},
					InsertOperation{Content: "Bar"},
					DeleteOperation{Content: "Foo"},
				},
			))
			Expect(new(ot.Transform).Insert("Fizz").Delete("Foo").Insert("Bar").Operations()).To(Equal(
				[]Operation{
					InsertOperation{Content: "FizzBar"},
					DeleteOperation{Content: "Foo"},
				},
			))
		})
	})

	Describe("Delete", func() {
		It("Does nothing when given an empty string", func() {
			Expect(new(ot.Transform).Delete("").Operations()).To(Equal([]Operation{}))
		})

		It("Adds to the base length", func() {
			Expect(new(ot.Transform).Delete("Foo").BaseLength).To(Equal(3))
			Expect(new(ot.Transform).Retain(3).Delete("Fizz").BaseLength).To(Equal(7))
		})

		It("Does not change the target length", func() {
			Expect(new(ot.Transform).Delete("Foo").TargetLength).To(Equal(0))
			Expect(new(ot.Transform).Retain(3).Delete("Fizz").TargetLength).To(Equal(3))
		})

		It("Saves the operation", func() {
			Expect(new(ot.Transform).Delete("Fizz").Operations()).To(Equal([]Operation{DeleteOperation{Content: "Fizz"}}))
		})

		It("Combines with previous Delete operations", func() {
			Expect(new(ot.Transform).Delete("Foo").Delete("Bar").Operations()).To(Equal([]Operation{DeleteOperation{Content: "FooBar"}}))
		})
	})
})
