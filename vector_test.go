package superstellar_test

import (
	. "superstellar/backend"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math"
)

var _ = Describe("Vector", func() {
	Describe("Length", func() {
		It("calculates length", func() {
			Expect((&Vector{X: 3, Y: 4}).Length()).To(Equal(5.0))
		})
	})

	Describe("Normalize", func() {
		var vector, expected *Vector

		Context("When vector is already normalized", func() {
			BeforeEach(func() {
				vector = &Vector{X: 1.0 / math.Sqrt2, Y: 1.0 / math.Sqrt2}
				expected = vector
			})

				It("does not change vector", func() {
				Expect(vector.Normalize()).To(Equal(expected))
			})
		})

		Context("When vector is not normalized", func() {
			BeforeEach(func() {
				vector = &Vector{X: 3, Y: 4}
				expected = &Vector{X: vector.X / 5, Y: vector.Y / 5}
			})

			It("normalizes vector", func() {
				Expect(vector.Normalize()).To(Equal(expected))
			})
		})
	})
})
