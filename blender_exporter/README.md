traytor blender exporter
========================

This is a python script for blender which exports the current scene into a
json file which can be used by the traytor.

# Usage:
Currently it's very stupid:

- open the python script in blender
- go to the bottom to the `json_to_file` function call and change the filename
- manually create a text block called `traytor_materials` and fill in your
  materials like so:

    "pretty_material": {
        "colour": [1, 0.2, 0],
        "type": "lambert"
    },
    "nice_lamp": {
        "colour": [1, 1, 1],
        "strength": 10000,
        "type": "emission"
    }

- run the script

This should save the entire scene to the json file you specified.

# Todo:
Make it a real exporter with material settings, nodes etc
