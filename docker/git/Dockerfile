FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV GITAMIX_FIXTURE_REPO=/opt/fixture/repo
ENV GITAMIX_FIXTURE_HASHES=/opt/fixture/repo/.hashes.fixture.env

RUN apt-get update \
	&& apt-get install -y --no-install-recommends bash ca-certificates git \
	&& rm -rf /var/lib/apt/lists/*

COPY docker/git/init_fake_git_tree.sh \
	/usr/local/bin/init_fake_git_tree

RUN chmod +x /usr/local/bin/init_fake_git_tree \
	&& /usr/local/bin/init_fake_git_tree "$GITAMIX_FIXTURE_REPO"

WORKDIR /opt/fixture/repo

CMD ["tail", "-f", "/dev/null"]