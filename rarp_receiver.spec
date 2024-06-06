Name:           rarp_receiver
Version:        1.0.0
Release:        1%{?dist}
Summary:        A simple RARP packet receiver written in Go

License:        MIT
URL:            https://github.com/artw/rarp_receiver 
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang

%description
A simple program that registers a receiver for RARP packets and prints the packet contents.

%prep
%setup -q

%build
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -ldflags="-linkmode external -extldflags '-Wl,--build-id'" -o rarp_receiver rarp_receiver.go

%install
install -d %{buildroot}%{_bindir}
install -m 0755 rarp_receiver %{buildroot}%{_bindir}/rarp_receiver

# Install systemd service file
install -d %{buildroot}%{_unitdir}
install -m 0644 rarp_receiver.service %{buildroot}%{_unitdir}/rarp_receiver.service

# Install default configuration file
install -d %{buildroot}%{_sysconfdir}/default
install -m 0644 rarp_receiver.default %{buildroot}%{_sysconfdir}/default/rarp_receiver

%files
%{_bindir}/rarp_receiver
%{_unitdir}/rarp_receiver.service
%{_sysconfdir}/default/rarp_receiver

%post
systemctl daemon-reload
systemctl enable rarp_receiver

%preun
if [ $1 -eq 0 ]; then
    systemctl disable rarp_receiver
fi

%postun
systemctl daemon-reload
systemctl restart rarp_receiver

%changelog
* Thu Jun 06 2024 Art Win <art@make.lv> - 1.0.0-1
- Initial package

