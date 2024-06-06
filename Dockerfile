FROM rockylinux:8

# Install required packages
RUN dnf install -y golang rpm-build && rpmdev-setuptree

# Set up build environment
WORKDIR /root

# Build the tarball for the source code
RUN 

# Build the RPM package
# RUN rpmbuild -ba rarp_receiver.spec

# The resulting RPM will be in /root/rpmbuild/RPMS/x86_64/

