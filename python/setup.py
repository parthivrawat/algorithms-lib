'''Setup configuration for algorithms-lib.'''

from setuptools import find_packages, setup

with open('README.md', 'r', encoding='utf-8') as fh:
    long_description = fh.read()

setup(
    name='algorithms-lib',
    version='1.0.0',
    author='Parthiv Rawat',
    license='MIT',
    description='A comprehensive, zero-dependency collection of common algorithms for Python',
    long_description=long_description,
    long_description_content_type='text/markdown',
    url='https://github.com/parthivrawat/algorithms-lib',
    packages=find_packages(),
    package_data={'algorithms_lib': ['py.typed']},
    include_package_data=True,
    classifiers=[
        'Development Status :: 4 - Beta',
        'Intended Audience :: Developers',
        'Topic :: Software Development :: Libraries :: Python Modules',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Programming Language :: Python :: 3.9',
        'Programming Language :: Python :: 3.10',
        'Programming Language :: Python :: 3.11',
        'Programming Language :: Python :: 3.12',
    ],
    python_requires='>=3.7',
    install_requires=[],
    extras_require={
        'dev': [
            'pytest>=7.0.0',
            'pytest-cov>=4.0.0',
        ],
    },
)
