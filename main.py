import bpy
import math
import numpy as np

#マテリアルを決める関数の定義
def material(name = 'material'):
    material_glass = bpy.data.materials.new(name)
    #ノードを使えるようにする
    material_glass.use_nodes = True
    p_BSDF = material_glass.node_tree.nodes["Principled BSDF"]
    #0→BaseColor/7→roughness(=粗さ)/15→transmission(=伝播)
    #default_value = (R, G, B, A)
    p_BSDF.inputs[0].default_value = np.random.rand(4)
    p_BSDF.inputs[7].default_value = 0
    p_BSDF.inputs[15].default_value = 1
    #オブジェクトにマテリアルの要素を追加する
    bpy.context.object.data.materials.append(material_glass)

#ランダムな値を3つ生成するための関数の定義
def get_random_location(min,max):
    return (max - min) * np.random.rand(3) + min

#random_number()関数が生成した値を使って立方体を回転、出現させる関数の定義
def generate_random_rotate_cube(min,max):
    random_location = get_random_location(min,max)
    bpy.ops.mesh.primitive_cube_add(location = random_location,rotation = math.pi * np.random.rand(3))

#複数の立方体をランダムな位置に回転、出現させる関数の定義
def generate_random_rotate_colorful_cubes(min,max,num):
    for i in range(0,num):
        generate_random_rotate_cube(min,max)
        material('Random')

generate_random_rotate_colorful_cubes(-70,70,500)
generate_random_rotate_colorful_cubes(-30,30,200)
generate_random_rotate_colorful_cubes(-10,10,50)

