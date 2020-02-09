#     Copyright (C) 2020, IrineSistiana
#
#     This file is part of mos-chinadns.
#
#     mosdns is free software: you can redistribute it and/or modify
#     it under the terms of the GNU General Public License as published by
#     the Free Software Foundation, either version 3 of the License, or
#     (at your option) any later version.
#
#     mosdns is distributed in the hope that it will be useful,
#     but WITHOUT ANY WARRANTY; without even the implied warranty of
#     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#     GNU General Public License for more details.
#
#     You should have received a copy of the GNU General Public License
#     along with this program.  If not, see <https://www.gnu.org/licenses/>.

# !/usr/bin/env python3
import subprocess
import sys
import os
import zipfile
import hashlib
import logging

project_name = 'mos-chinadns'
main_path = './'
CGO_ENABLED = '0'

do_upx = True
do_zip = True

# [(env : value),(env : value)]
env_combinations = [[['GOOS', 'darwin'], ['GOARCH', 'amd64']],

                    # [['GOOS', 'linux'], ['GOARCH', '386']],
                    [['GOOS', 'linux'], ['GOARCH', 'amd64']],

                    # [['GOOS', 'linux'], ['GOARCH', 'arm'], ['GOARM', '7']],
                    # [['GOOS', 'linux'], ['GOARCH', 'arm64']],

                    # [['GOOS', 'linux'], ['GOARCH', 'mips'], ['GOMIPS', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips'], ['GOMIPS', 'softfloat']],
                    #[['GOOS', 'linux'], ['GOARCH', 'mipsle'], ['GOMIPS', 'hardfloat']],
                    [['GOOS', 'linux'], ['GOARCH', 'mipsle'], ['GOMIPS', 'softfloat']],

                    # [['GOOS', 'linux'], ['GOARCH', 'mips64'],
                    #  ['GOMIPS64', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips64'],
                    #  ['GOMIPS64', 'softfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips64le'],
                    #  ['GOMIPS64', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips64le'],
                    #  ['GOMIPS64', 'softfloat']],

                    # [['GOOS', 'freebsd'], ['GOARCH', '386']],
                    # [['GOOS', 'freebsd'], ['GOARCH', 'amd64']],

                    # [['GOOS', 'windows'], ['GOARCH', '386']],
                    [['GOOS', 'windows'], ['GOARCH', 'amd64']],
                    ]



def go_build():
    logging.warning(f'building {project_name}')

    global env_combinations
    if len(sys.argv) == 2 and sys.argv[1].isdigit():
        index = int(sys.argv[1])
        env_combinations = [env_combinations[index]]

    for env in env_combinations:
        os_env = os.environ.copy()  # new env

        zip_name = project_name
        for pairs in env:
            os_env[pairs[0]] = pairs[1]  # add env
            zip_name = zip_name + '-' + pairs[1]
        os_env['CGO_ENABLED'] = CGO_ENABLED

        binary_name = project_name + ('.exe' if os_env['GOOS'] == 'windows' else '')

        logging.warning(f'building {zip_name}')
        try:
            go_build_ldflags = '-s -w'
            go_build_other_args = ''

            buildCmd = f'go build -ldflags "{go_build_ldflags}" {go_build_other_args} -o {binary_name} {main_path}'
            subprocess.check_call(buildCmd, shell=True, env=os_env)
            upx_path = 'upx'
            if do_upx:
                try:
                    subprocess.check_call(f'{upx_path} -9 -q {binary_name}', shell=True, stderr=subprocess.DEVNULL,
                                          stdout=subprocess.DEVNULL)
                except subprocess.CalledProcessError as e:
                    logging.error(f'upx failed: {e.args}')

            zip_includes = []
            if do_zip:
                with zipfile.ZipFile(f'{zip_name}.zip', mode='w', compression=zipfile.ZIP_DEFLATED,
                                     compresslevel=1) as zf:
                    for path in zip_includes:
                        if os.path.isdir(path):  # path is a dir
                            for dirpath, dirnames, filenames in os.walk('common'):
                                for filename in filenames:
                                    zf.write(os.path.join(dirpath, filename), filename)
                        elif os.path.isfile(path):
                            zf.write(path, os.path.basename(path))
                        else:
                            logging.warning(f'path is not a dir nor a file: {path}')

                    zf.write(binary_name)

        except subprocess.CalledProcessError as e:
            logging.error(f'build {zip_name} failed: {e.args}')
        except Exception:
            logging.exception('unknown err')


if __name__ == '__main__':
    go_build()
