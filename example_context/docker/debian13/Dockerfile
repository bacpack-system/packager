FROM debian:13.0

USER root
RUN echo root:1234 | chpasswd

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    binutils build-essential flex git libssl-dev lsb-release make pkg-config \
    openssh-server patchelf wget libmpfr-dev libmpc-dev libgmp-dev && \
    rm -rf /var/lib/apt/lists/*

RUN wget "https://github.com/Kitware/CMake/releases/download/v3.25.2/cmake-3.25.2-linux-x86_64.sh" -O cmake.sh && \
    chmod +x cmake.sh && \
    ./cmake.sh --skip-license --prefix=/usr/local && \
    rm ./cmake.sh

RUN apt-get update && \
    apt-get purge -y \
      wget && \
      rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/cmakelib/cmakelib.git /cmakelib
RUN echo "export CMLIB_DIR=/cmakelib" >> /root/.bashrc

# Install Fleet Protocol's CMCONF system
COPY --from=package-context config/CMCONF_FLEET_PROTOCOLConfig.cmake /etc
ENV CMLIB_DIR=/cmakelib
RUN cmake -DCMCONF_INSTALL_AS_SYMLINK=ON -P /etc/CMCONF_FLEET_PROTOCOLConfig.cmake

RUN sed -ri 's/#?PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN mkdir -p /run/sshd

ENTRYPOINT ["/usr/sbin/sshd", "-D", "-o", "ListenAddress=0.0.0.0"]
