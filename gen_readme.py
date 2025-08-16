import os
import subprocess

readme_contents = open('README_template.md', 'r').read()

def _replace_section(contents, section_name, package_name):
    section_content = subprocess.check_output(
        ['go', 'doc', '-all', package_name]
    ).decode("utf-8")
    section_content = _markdownize(section_content)
    return contents.replace(f'$${section_name}$$', section_content)


readme_contents = _replace_section(readme_contents, 'FUNCTIONAL_DOCS', './functional')
readme_contents = _replace_section(readme_contents, 'SETS_DOCS', './sets')

with open('README.md', 'w') as readme_file:
    readme_file.write(readme_contents)