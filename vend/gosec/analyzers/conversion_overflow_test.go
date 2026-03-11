package analyzers

import (
	"go/types"
	"math"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetIntTypeInfo", func() {
	Context("with valid input", func() {
		DescribeTable("should correctly parse and calculate bounds for",
			func(kind types.BasicKind, expectedSigned bool, expectedSize int, expectedMin int64, expectedMax uint64) {
				// Use the standard shared basic types directly
				basicType := types.Typ[kind]

				result, err := GetIntTypeInfo(basicType)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Signed).To(Equal(expectedSigned))
				Expect(result.Size).To(Equal(expectedSize))
				Expect(result.Min).To(Equal(expectedMin))
				Expect(result.Max).To(Equal(expectedMax))
			},
			Entry("uint8", types.Uint8, false, 8, int64(0), uint64(math.MaxUint8)),
			Entry("int8", types.Int8, true, 8, int64(math.MinInt8), uint64(math.MaxInt8)),
			Entry("uint16", types.Uint16, false, 16, int64(0), uint64(math.MaxUint16)),
			Entry("int16", types.Int16, true, 16, int64(math.MinInt16), uint64(math.MaxInt16)),
			Entry("uint32", types.Uint32, false, 32, int64(0), uint64(math.MaxUint32)),
			Entry("int32", types.Int32, true, 32, int64(math.MinInt32), uint64(math.MaxInt32)),
			Entry("uint64", types.Uint64, false, 64, int64(0), uint64(math.MaxUint64)),
			Entry("int64", types.Int64, true, 64, int64(math.MinInt64), uint64(math.MaxInt64)),
		)

		It("should use system's int size for 'int' and 'uint'", func() {
			intResult, err := GetIntTypeInfo(types.Typ[types.Int])
			Expect(err).NotTo(HaveOccurred())
			Expect(intResult.Size).To(Equal(strconv.IntSize))

			uintResult, err := GetIntTypeInfo(types.Typ[types.Uint])
			Expect(err).NotTo(HaveOccurred())
			Expect(uintResult.Size).To(Equal(strconv.IntSize))
		})
	})

	Context("with invalid input", func() {
		It("should return error for non-basic types", func() {
			_, err := GetIntTypeInfo(types.NewSlice(types.Typ[types.Int]))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a basic type"))
		})

		It("should return error for non-integer basic types", func() {
			_, err := GetIntTypeInfo(types.Typ[types.Float64])
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported basic type"))
		})
	})
})

// Helper to simulate isIntOverflow logic using GetIntTypeInfo
func checkOverflow(srcKind, dstKind types.BasicKind) bool {
	srcInfo, err := GetIntTypeInfo(types.Typ[srcKind])
	if err != nil {
		return false
	}
	dstInfo, err := GetIntTypeInfo(types.Typ[dstKind])
	if err != nil {
		return false
	}
	return hasOverflow(srcInfo, dstInfo)
}

var _ = Describe("Overflow Logic (simulated)", func() {
	DescribeTable("should correctly identify overflow scenarios",
		func(src types.BasicKind, dst types.BasicKind, expectedOverflow bool) {
			Expect(checkOverflow(src, dst)).To(Equal(expectedOverflow))
		},
		// Unsigned to Signed conversions
		Entry("uint8 to int8", types.Uint8, types.Int8, true),
		Entry("uint8 to int16", types.Uint8, types.Int16, false),
		Entry("uint8 to int32", types.Uint8, types.Int32, false),
		Entry("uint8 to int64", types.Uint8, types.Int64, false),
		Entry("uint16 to int8", types.Uint16, types.Int8, true),
		Entry("uint16 to int16", types.Uint16, types.Int16, true),
		Entry("uint16 to int32", types.Uint16, types.Int32, false),
		Entry("uint16 to int64", types.Uint16, types.Int64, false),
		Entry("uint32 to int8", types.Uint32, types.Int8, true),
		Entry("uint32 to int16", types.Uint32, types.Int16, true),
		Entry("uint32 to int32", types.Uint32, types.Int32, true),
		Entry("uint32 to int64", types.Uint32, types.Int64, false),
		Entry("uint64 to int8", types.Uint64, types.Int8, true),
		Entry("uint64 to int16", types.Uint64, types.Int16, true),
		Entry("uint64 to int32", types.Uint64, types.Int32, true),
		Entry("uint64 to int64", types.Uint64, types.Int64, true),

		// Unsigned to Unsigned conversions
		Entry("uint8 to uint16", types.Uint8, types.Uint16, false),
		Entry("uint8 to uint32", types.Uint8, types.Uint32, false),
		Entry("uint8 to uint64", types.Uint8, types.Uint64, false),
		Entry("uint16 to uint8", types.Uint16, types.Uint8, true),
		Entry("uint16 to uint32", types.Uint16, types.Uint32, false),
		Entry("uint16 to uint64", types.Uint16, types.Uint64, false),
		Entry("uint32 to uint8", types.Uint32, types.Uint8, true),
		Entry("uint32 to uint16", types.Uint32, types.Uint16, true),
		Entry("uint32 to uint64", types.Uint32, types.Uint64, false),
		Entry("uint64 to uint8", types.Uint64, types.Uint8, true),
		Entry("uint64 to uint16", types.Uint64, types.Uint16, true),
		Entry("uint64 to uint32", types.Uint64, types.Uint32, true),

		// Signed to Unsigned conversions
		Entry("int8 to uint8", types.Int8, types.Uint8, true),
		Entry("int8 to uint16", types.Int8, types.Uint16, true),
		Entry("int8 to uint32", types.Int8, types.Uint32, true),
		Entry("int8 to uint64", types.Int8, types.Uint64, true),
		Entry("int16 to uint8", types.Int16, types.Uint8, true),
		Entry("int16 to uint16", types.Int16, types.Uint16, true),
		Entry("int16 to uint32", types.Int16, types.Uint32, true),
		Entry("int16 to uint64", types.Int16, types.Uint64, true),
		Entry("int32 to uint8", types.Int32, types.Uint8, true),
		Entry("int32 to uint16", types.Int32, types.Uint16, true),
		Entry("int32 to uint32", types.Int32, types.Uint32, true),
		Entry("int32 to uint64", types.Int32, types.Uint64, true),
		Entry("int64 to uint8", types.Int64, types.Uint8, true),
		Entry("int64 to uint16", types.Int64, types.Uint16, true),
		Entry("int64 to uint32", types.Int64, types.Uint32, true),
		Entry("int64 to uint64", types.Int64, types.Uint64, true),

		// Signed to Signed conversions
		Entry("int8 to int16", types.Int8, types.Int16, false),
		Entry("int8 to int32", types.Int8, types.Int32, false),
		Entry("int8 to int64", types.Int8, types.Int64, false),
		Entry("int16 to int8", types.Int16, types.Int8, true),
		Entry("int16 to int32", types.Int16, types.Int32, false),
		Entry("int16 to int64", types.Int16, types.Int64, false),
		Entry("int32 to int8", types.Int32, types.Int8, true),
		Entry("int32 to int16", types.Int32, types.Int16, true),
		Entry("int32 to int64", types.Int32, types.Int64, false),
		Entry("int64 to int8", types.Int64, types.Int8, true),
		Entry("int64 to int16", types.Int64, types.Int16, true),
		Entry("int64 to int32", types.Int64, types.Int32, true),

		// Same type conversions (should never overflow)
		Entry("uint8 to uint8", types.Uint8, types.Uint8, false),
		Entry("uint16 to uint16", types.Uint16, types.Uint16, false),
		Entry("uint32 to uint32", types.Uint32, types.Uint32, false),
		Entry("uint64 to uint64", types.Uint64, types.Uint64, false),
		Entry("int8 to int8", types.Int8, types.Int8, false),
		Entry("int16 to int16", types.Int16, types.Int16, false),
		Entry("int32 to int32", types.Int32, types.Int32, false),
		Entry("int64 to int64", types.Int64, types.Int64, false),
	)
})
