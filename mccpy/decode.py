from array import array
import struct
import math
import sys

def encode_traytor_hdr(format, image):
    size = array('h', image.size)
    pixels = array('f', image.pixels)   # note: this is huge!
    
    if format == 'traytor_srgb':
        for i in range(len(pixels)):
            pixels[i] = srgb_to_linear(pixels[i])

    return size.tobytes() + pixels.tobytes()

def decode_traytor_srgb(array_of_bytes):
    size = struct.unpack('hh', array_of_bytes[0:4])
    pixel_bytes = array_of_bytes[4:]
    
    pixels = array('f')
    pixels.frombytes(pixel_bytes)

    return pixels.tolist()

def read_binary_file(filen):
    with open(filen, "rb") as f:
        bytes = f.read()
    return bytes

def compare(file1, file2):
    first = decode_traytor_srgb(read_binary_file(file1))
    second = decode_traytor_srgb(read_binary_file(file2))
   
    return math.sqrt(sum([(a - b) ** 2 for a, b in zip(first, second)]))

if __name__ == '__main__':
    if len(sys.argv) == 3:
        f1 = sys.argv[1]
        f2 = sys.argv[2]
        print(compare(f1, f2))
    else:
        print("Wrong number of arguments")