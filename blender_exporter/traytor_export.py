import bpy
import json
import bmesh
import mathutils
import gzip
import base64
import os
from array import array
from tempfile import mkstemp

def make_material_index(mesh, face_material_index):
    material = mesh.materials[face_material_index]
    all_materials = bpy.data.materials
    return [i for i in range(len(all_materials)) if material.name == all_materials[i].name][0]

def make_material(mesh, material):
    return {
        'name': material.name
    }

def make_vertex(mesh, index):
    vertex = mesh.vertices[mesh.loops[index].vertex_index]
    data = {
        'coordinates': list(vertex.co),
        'normal': list(vertex.normal)
    }
    uv_layer = mesh.uv_layers.active
    if uv_layer:
        data['uv'] = list(uv_layer.data[index].uv)

    return data

def make_face(mesh, face, vertex_index_offset):
    face_vertices = face.loop_indices
    data = {
        'vertices': [vertex + vertex_index_offset for vertex in face_vertices],
        'material': make_material_index(mesh, face.material_index)
    }
    
    if not face.use_smooth:
        data['normal'] = list(face.normal)
    
    return data

def make_camera(camera, scene):
    transformation = camera.matrix_world
    
    frame = [list(transformation * point) for point in camera.data.view_frame(scene)]
    
    data = {
        'type': 'pinhole',
        'top_right': frame[0],
        'bottom_right': frame[1],
        'bottom_left': frame[2],
        'top_left': frame[3],
        'focus': list(transformation * mathutils.Vector([0, 0, 0]))
    }
    
    return data

def get_materials(mesh):
    return [make_material(mesh, material) for material in mesh.materials]

def get_vertices(mesh):
    return [make_vertex(mesh, index) for index in range(len(mesh.loops))]

def get_faces(mesh, vertex_index_offset):
    return [make_face(mesh, face, vertex_index_offset) for face in mesh.polygons]

def srgb_to_linear(value):
    if value < 0.04045:
        return value / 12.92
    return pow((value + 0.055)/1.055, 2.4)

def encode_traytor_hdr(format, image):
    size = array('h', image.size)
    pixels = array('f', image.pixels)   # note: this is huge!
    
    if format == 'traytor_srgb':
        for i in range(len(pixels)):
            pixels[i] = srgb_to_linear(pixels[i])

    return size.tobytes() + pixels.tobytes()

def image_data(image_name, format = None):
    image = bpy.data.images[image_name]
    if format:
        if format == 'traytor_hdr' or format == 'traytor_srgb':
            return 'traytor_hdr', base64.b64encode(
                encode_traytor_hdr(format, image)
            ).decode('utf-8')
        else:
            image.file_format = format
    else:
        format = image.file_format
    
    f, filename = mkstemp()
    os.close(f)
    
    try:
        image.filepath_raw = filename
        image.save()
        
        with open(filename, "rb") as f:
            return format, base64.b64encode(f.read()).decode('utf-8')
    finally:
        os.remove(filename)
    
def walk_materials(data):
    if not isinstance(data, dict):
        return
    
    if data.get('type') == 'image_texture' and 'image' in data:
        data['format'], data['data'] = image_data(
            data['image'], data.get('format')
        )
        
    for _, item in data.items():
        walk_materials(item)

def expand_materials(materials):
    material_defs = bpy.data.texts['traytor_materials'].as_string()
    material_data = json.loads(material_defs)
    walk_materials(material_data)
    return [material_data[material['name']] for material in materials]

def triangulate(mesh):
    bm = bmesh.new()
    bm.from_mesh(mesh)
    bmesh.ops.triangulate(bm, faces=bm.faces)
    bm.to_mesh(mesh)
    bm.free()
    
def get_scene(scene):
    vertices = []
    faces = []
    
    for obj in scene.objects:
        if obj.type == 'MESH':
            transformation = obj.matrix_world

            mesh = obj.to_mesh(scene, apply_modifiers=True, settings='RENDER')
            try:
                triangulate(mesh)
                mesh.transform(transformation)
                mesh.calc_normals()
                
                faces += get_faces(mesh, len(vertices))
                vertices += get_vertices(mesh)
            finally:                
                bpy.data.meshes.remove(mesh)
    
    data = json.loads(bpy.data.texts['traytor_settings'].as_string())
    
    data['mesh'] = {
        'vertices': vertices,
        'faces': faces,
    }
    data['materials'] = expand_materials(
        [make_material(mesh, m) for m in bpy.data.materials]
    )
    
    if scene.camera:
        data['camera'] = make_camera(scene.camera, scene)
        
    return data
           

def json_to_file(scene, file):
    with open(file, 'w') as f:
        json.dump(get_scene(scene), f, sort_keys=True, indent=4)

def jsongz_to_file(scene, file):
    with gzip.open(file, 'wt') as f:
        json.dump(get_scene(scene), f, separators=(',', ':'))

      
json_to_file(bpy.context.scene, '/tmp/scene.json')
jsongz_to_file(bpy.context.scene, '/tmp/scene.json.gz')
