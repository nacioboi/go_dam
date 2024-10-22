# Things i learned.

- Using multiplier can work so long as we store the remainder as well. This is because not all values divide evenly by the range of multipliers 1-16.

- Instead it might be easier to just use more bits to store the difference.
  - And thats what i did. I used 30 bits instead of 16.

- For random data, the compressed array is only slightly better but with a gigantic performance hit.

## It is worth it to explore huffman bit compression.

1. Have chunks of uint64s and compress them using a global huffman tree.
2. Then we can have a workspace arr where we can store the decompressed data.
3. If we get this just right, we should be able to compress almost any data type but not completely random data.
