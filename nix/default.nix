{ system ? builtins.currentSystem }:
let
  pkgs = (import ./nixpkgs.nix {
    overlays = [
      (final: pkg: {
        pcre = (static pkg.pcre).overrideAttrs (x: {
          configureFlags = x.configureFlags ++ [
            "--enable-static"
          ];
        });
      })
    ];
    config = {
      packageOverrides = pkg: {
        autogen = (static pkg.autogen);
        e2fsprogs = (static pkg.e2fsprogs);
        gnupg = (static pkg.gnupg);
        gpgme = (static pkg.gpgme);
        libassuan = (static pkg.libassuan);
        libgpgerror = (static pkg.libgpgerror);
        libseccomp = (static pkg.libseccomp);
        libuv = (static pkg.libuv);
        glib = (static pkg.glib).overrideAttrs (x: {
          outputs = [ "bin" "out" "dev" ];
          mesonFlags = [
            "-Ddefault_library=static"
            "-Ddevbindir=${placeholder ''dev''}/bin"
            "-Dgtk_doc=false"
            "-Dnls=disabled"
          ];
          postInstall = ''
            moveToOutput "share/glib-2.0" "$dev"
            substituteInPlace "$dev/bin/gdbus-codegen" --replace "$out" "$dev"
            sed -i "$dev/bin/glib-gettextize" -e "s|^gettext_dir=.*|gettext_dir=$dev/share/glib-2.0/gettext|"
            sed '1i#line 1 "${x.pname}-${x.version}/include/glib-2.0/gobject/gobjectnotifyqueue.c"' \
              -i "$dev"/include/glib-2.0/gobject/gobjectnotifyqueue.c
          '';
        });
        gnutls = (static pkg.gnutls).overrideAttrs (x: {
          configureFlags = (x.configureFlags or [ ]) ++ [
            "--disable-non-suiteb-curves"
            "--disable-openssl-compatibility"
            "--disable-rpath"
            "--enable-local-libopts"
            "--without-p11-kit"
          ];
        });
        systemd = (static pkg.systemd).overrideAttrs (x: {
          outputs = [ "out" "dev" ];
          mesonFlags = x.mesonFlags ++ [
            "-Dstatic-libsystemd=true"
          ];
        });
      };
    };
  });

  static = pkg: pkg.overrideAttrs (x: {
    doCheck = false;
    configureFlags = (x.configureFlags or [ ]) ++ [
      "--without-shared"
      "--disable-shared"
    ];
    dontDisableStatic = true;
    enableSharedExecutables = false;
    enableStatic = true;
  });

  self = with pkgs; buildGoModule rec {
    name = "skopeo";
    src = ./..;
    vendorSha256 = null;
    doCheck = false;
    enableParallelBuilding = true;
    outputs = [ "out" ];
    nativeBuildInputs = [ bash go-md2man installShellFiles makeWrapper pcre pkg-config which ];
    buildInputs = [ glibc glibc.static glib gpgme libassuan libgpgerror libseccomp ];
    prePatch = ''
      export CFLAGS='-static -pthread'
      export LDFLAGS='-s -w -static-libgcc -static'
      export EXTRA_LDFLAGS='-s -w -linkmode external -extldflags "-static -lm"'
      export BUILDTAGS='static netgo osusergo exclude_graphdriver_btrfs exclude_graphdriver_devicemapper'
    '';
    buildPhase = ''
      patchShebangs .
      make bin/skopeo
    '';
    installPhase = ''
      install -Dm755 bin/skopeo $out/bin/skopeo
    '';
  };
in
self
