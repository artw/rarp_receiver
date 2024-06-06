FROM rockylinux:8

# Install required packages
RUN dnf install -y golang rpm-build && mkdir -p rpmbuild/SOURCES

# Set up build environment
WORKDIR /root

# Build the tarball for the source code
RUN 

# Build the RPM package
ENTRYPOINT ["rpmbuild","-ba"]
# RUN rpmbuild -ba myrpm.spec

# The resulting RPM will be in /root/rpmbuild/RPMS/x86_64/

