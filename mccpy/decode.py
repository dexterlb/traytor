from array import array
import struct

def encode_traytor_hdr(format, image):
    size = array('h', image.size)
    pixels = array('f', image.pixels)   # note: this is huge!
    
    if format == 'traytor_srgb':
        for i in range(len(pixels)):
            pixels[i] = srgb_to_linear(pixels[i])

    return size.tobytes() + pixels.tobytes()

def decode_traytor_srgb(array_of_bytes):
	size = struct.unpack('h', array_of_bytes[0:2])[0], struct.unpack('h', array_of_bytes[2:4])[0]
	pixel_bytes = array_of_bytes[4:]
	pixels = []
	for i in range(0, len(pixel_bytes), 4):
		pixels.append(struct.unpack('f', pixel_bytes[i: i+4])[0])
	return pixels
