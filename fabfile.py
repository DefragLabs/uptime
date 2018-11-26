import os

from fabric.api import task, run, abort, settings, local
from fabric.context_managers import cd
from fabric.contrib import files

CODE_DIRECTORY = "/home/ubuntu/code"
BACKEND_REPO_URL = "https://github.com/DefragLabs/uptime.git"
FRONTEND_REPO_URL = "https://github.com/DefragLabs/uptime-ui.git"
SLACK_WEBHOOK_URL = ""

DEV_HOST_STRING = "ubuntu@18.235.105.251"

KEY_NAME = "uptime.pem"
DEV_KEY_NAME = "uptime.pem"
KEY_DIRECTORY = "~/.ssh/"


def extract_path_from_url(url):
    """
    Possible formats.

    https://github.com/DefragLabs/uptime.git

    The directory would be `uptime`
    """
    resource = url.split("/")[-1]
    resource_without_extension = resource.split(".")[0]
    return os.path.join(CODE_DIRECTORY, resource_without_extension)


def clone_repo(repo_url):
    repo_path = extract_path_from_url(repo_url)

    if not files.exists(repo_path):
        run('git clone {}'.format(repo_url))
    else:
        print("Repo already exists.")


def checkout_and_update_branch(branch):
    run('git fetch')
    run('git checkout {}'.format(branch))
    run('git pull --rebase origin {}'.format(branch))


def notify_on_slack(msg):
    pass


def deploy_frontend(branch, env="dev"):
    clone_repo(FRONTEND_REPO_URL)

    frontend_repo_path = extract_path_from_url(FRONTEND_REPO_URL)

    with cd(frontend_repo_path):
        checkout_and_update_branch(branch)

        print("Running development setup.")
        if env == "prod":
            run('docker-compose -f docker-compose-prod.yml down')
            run('docker-compose -f docker-compose-prod.yml up -d --build')
        elif env == "dev":
            run('docker-compose -f docker-compose-dev.yml down')
            run('docker-compose -f docker-compose-dev.yml up -d --build')
        elif env == "uat":
            run('docker-compose -f docker-compose-uat.yml down')
            run('docker-compose -f docker-compose-uat.yml up -d --build')
        else:
            abort("Invalid env")
        short_commit_hash = run('git rev-parse --short HEAD')

        msg = "Deployed frontend version: {} to {}".format(short_commit_hash, env)
        notify_on_slack(msg)


def deploy_backend(branch, env=None):
    clone_repo(BACKEND_REPO_URL)

    repo_path = extract_path_from_url(BACKEND_REPO_URL)
    with cd(repo_path):
        checkout_and_update_branch(branch)
        print("Restarting backend server.")
        if env == "prod":
            run('docker-compose -f docker-compose-prod.yml down')
            run('docker-compose -f docker-compose-prod.yml up -d --build')
        elif env == "dev":
            run('docker-compose -f docker-compose-dev.yml down')
            run('docker-compose -f docker-compose-dev.yml up -d --build')
        elif env == "uat":
            run('docker-compose -f docker-compose-uat.yml down')
            run('docker-compose -f docker-compose-uat.yml up -d --build')
        else:
            abort("Invalid env")
        short_commit_hash = run('git rev-parse --short HEAD')
        msg = "Deployed backend version: {} to {}".format(short_commit_hash, env)
        notify_on_slack(msg)


@task
def deploy_dev(app=None, branch='master'):
    key_name = DEV_KEY_NAME
    key_path = KEY_DIRECTORY + key_name
    if not os.path.isfile(os.path.expanduser(key_path)):
        abort("Unable to find ssh key at: %s" % key_path)

    with settings(host_string=DEV_HOST_STRING, key_filename=key_path, command_timeout=1000):
        run('mkdir -p {}'.format(CODE_DIRECTORY))

        # to solve no space left on device issue.
        run('docker image prune -f && docker container prune -f')
        with cd(CODE_DIRECTORY):
            if app is None:
                deploy_frontend(branch, env="dev")
                deploy_backend(branch, env="dev")
            elif app.lower() == 'backend':
                deploy_backend(branch, env="dev")
            elif app.lower() == 'frontend':
                deploy_frontend(branch, env="dev")
            else:
                abort("Unrecognized app {}".format(app))
