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

def compare(first, second):
    return math.sqrt(sum([(a - b) ** 2 for a, b in zip(first, second)]))

def norm(original):
    return math.sqrt(sum(a ** 2 for a in original))

if __name__ == '__main__':
    if len(sys.argv) == 2:
        original_file = "/home/do/go/src/github.com/DexterLB/traytor/mccpy/new/normal_3000.th"
        file_name = sys.argv[1]

        original = decode_traytor_srgb(read_binary_file(original_file))
        second = decode_traytor_srgb(read_binary_file(file_name))

        perr = compare(original, second)
        print("Dist: %.3f; Percent error: %.3f%%" % (perr, perr/norm(original) * 100))
    else:
        print("Wrong number of arguments")