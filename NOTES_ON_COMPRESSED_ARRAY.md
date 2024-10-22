# Things i learned.

- Using multiplier can work so long as we store the remainder as well. This is because not all values divide evenly by the range of multipliers 1-16.

- Instead it might be easier to just use more bits to store the difference.
  - And thats what i did. I used 30 bits instead of 16.

- For random data, the compressed array is only slightly better but with a gigantic performance hit.

## It is worth it to explore huffman bit compression.

1. Have chunks of uint64s and compress them using a global huffman tree.
2. Then we can have a workspace arr where we can store the decompressed data.
3. If we get this just right, we should be able to compress almost any data type but not completely random data.

## Also explore

Golomb-Rice coding is essentially about efficiently storing values in a bit-stream by encoding numbers with two parts:

- Quotient (unary encoding): This is the number of times the value is divisible by a certain base (usually a power of 2, for Rice coding). This part is encoded as a series of 1s followed by a 0.

-Remainder (binary encoding): This is the remainder when dividing by the base, stored in a fixed number of bits.
For small numbers, the encoding requires fewer bits than the original value. The key is selecting a base (often 2^k) that best matches the typical range of your data to minimize bit usage.

In simple terms, it "compresses" by breaking down the value into chunks, where the smaller chunk sizes get the fewest bits in the stream, but it works best when your data values are small or predictable.
