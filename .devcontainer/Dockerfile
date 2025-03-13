FROM debian:12-slim

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Configure apt and install packages
RUN apt-get update

# Install git, process tools, lsb-release (common in install instructions for CLIs) and other dependencies
RUN apt-get -y install git procps lsb-release

# Install Development Tools
RUN apt-get -y install vim zsh curl wget unzip sqlite3 gcc make build-essential


RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

RUN echo "export PATH=/usr/local/go/bin:${PATH}" >> /root/.zshrc

RUN chsh -s /bin/zsh root

RUN git config --global --add safe.directory /workspaces/powered-storage

ENTRYPOINT [ "/bin/zsh" ]