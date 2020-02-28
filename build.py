
# !/usr/bin/env python3
import subprocess
import sys
import os
import zipfile
import hashlib
import logging

import urllib3

project_name = 'mos-chinadns'

# [(env : value),(env : value)]
env_combinations = [[['GOOS', 'darwin'], ['GOARCH', 'amd64']],

                    # [['GOOS', 'linux'], ['GOARCH', '386']],
                    [['GOOS', 'linux'], ['GOARCH', 'amd64']],

                    # [['GOOS', 'linux'], ['GOARCH', 'arm'], ['GOARM', '7']],
                    # [['GOOS', 'linux'], ['GOARCH', 'arm64']],

                    # [['GOOS', 'linux'], ['GOARCH', 'mips'], ['GOMIPS', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips'], ['GOMIPS', 'softfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mipsle'],
                    #     ['GOMIPS', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mipsle'],
                    #     ['GOMIPS', 'softfloat']],

                    # [['GOOS', 'linux'], ['GOARCH', 'mips64'],
                    #  ['GOMIPS64', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips64'],
                    #  ['GOMIPS64', 'softfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips64le'],
                    #  ['GOMIPS64', 'hardfloat']],
                    # [['GOOS', 'linux'], ['GOARCH', 'mips64le'],
                    #  ['GOMIPS64', 'softfloat']],

                    # [['GOOS', 'freebsd'], ['GOARCH', '386']],
                    [['GOOS', 'freebsd'], ['GOARCH', 'amd64']],

                    # [['GOOS', 'windows'], ['GOARCH', '386']],
                    [['GOOS', 'windows'], ['GOARCH', 'amd64']],
                    ]

def init_release_resources():
    subprocess.check_call(f'go run ./ -gen config-example.json', shell=True, env=os.environ.copy())
    url = 'https://raw.githubusercontent.com/LisonFan/china_ip_list/master/china_ipv4_ipv6_list' 
    print(f'downloading chn list') 
    http = urllib3.PoolManager()
    response = http.request('GET', url)
    with open('chn.list', 'wb') as f:
        f.write(response.data)
    response.release_conn()
    print(f'chn list retrieved') 

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

        suffix = '.exe' if os_env['GOOS'] == 'windows' else ''

        logging.warning(f'building {zip_name}')
        try:
            subprocess.check_call(
                f'go build -ldflags "-s -w"  -o {project_name}{suffix}', shell=True, env=os_env)
            try:
                subprocess.check_call(f'upx -9 -q {project_name}{suffix}', shell=True, stderr=subprocess.DEVNULL,
                                      stdout=subprocess.DEVNULL)
            except subprocess.CalledProcessError as e:
                logging.error(f'upx failed: {e.args}')

            with zipfile.ZipFile(f'{zip_name}.zip', mode='w', compression=zipfile.ZIP_DEFLATED,
                                 compresslevel=5  ) as zf:
                zf.write(f'{project_name}{suffix}')
                zf.write('README.md')
                zf.write('config-example.json')
                zf.write('chn.list')
                zf.write('LICENSE')

        except subprocess.CalledProcessError as e:
            logging.error(f'build {zip_name} failed: {e.args}')
        except Exception:
            logging.exception('unknown err')


if __name__ == '__main__':
    init_release_resources()
    go_build()
