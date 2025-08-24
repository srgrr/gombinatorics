import os
import subprocess
import sys

README_TEMPLATE = "README_template.md"
SAMPLE_FILE = "examples/range_filter_map/range_filter_map.go"
DECLARATIVE_FILE = "examples/declarative_example/declarative_abbr.txt"
FLASHY_FILE = "examples/range_filter_map/flashy.txt"
README_FILE = "README.md"
VERSION_TAG = "<<VERSION>>"
SAMPLE_TAG = "<<SAMPLE>>"
DECLARATIVE_TAG = "<<DECLARATIVE>>"
FLASHY_TAG = "<<FLASHY_SAMPLE>>"

CURRENT_VERSION = sys.argv[1]


with open(README_TEMPLATE, 'r') as template_file:
    content = template_file.read()

content = content.replace(VERSION_TAG, CURRENT_VERSION).replace(SAMPLE_TAG, open(SAMPLE_FILE, 'r').read()).replace(DECLARATIVE_TAG, open(DECLARATIVE_FILE, 'r').read()).replace(FLASHY_TAG, open(FLASHY_FILE, 'r').read())

with open(README_FILE, 'w') as readme_file:
    readme_file.write(content)

