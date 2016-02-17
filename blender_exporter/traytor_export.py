import bpy
import json
import bmesh
import mathutils
import gzip

def make_material_index(mesh, face_material_index):
    material = mesh.materials[face_material_index]
    all_materials = bpy.data.materials
    return [i for i in range(len(all_materials)) if material.name == all_materials[i].name][0]

def make_material(mesh, material):
    return {
        'name': material.name
    }

def make_vertex(mesh, index):
    uv_layer = mesh.uv_layers.active
    vertex = mesh.vertices[index]
    data = {
        'coordinates': list(vertex.co),
        'normal': list(vertex.normal)
    }
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
    return [make_vertex(mesh, index) for index in range(len(mesh.vertices))]

def get_faces(mesh, vertex_index_offset):
    return [make_face(mesh, face, vertex_index_offset) for face in mesh.polygons]

def expand_materials(materials):
    material_defs = bpy.data.texts['traytor_materials'].as_string()
    material_data = json.loads(material_defs)
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
    
    data = {
        'mesh': {
            'vertices': vertices,
            'faces': faces,
        },
        'materials': expand_materials([make_material(mesh, m) for m in bpy.data.materials])
    }
    
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