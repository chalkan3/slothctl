lvim_install_prerequisites:
  pkg.installed:
    - names:
      - git
      - neovim
      - make
      - python-pip
      - npm
      - nodejs
      - cargo
      - ripgrep
      - curl # Needed for the installer script

lvim_install:
  cmd.run:
    - name: LV_BRANCH='release-1.4/neovim-0.9' bash <(curl -s https://raw.githubusercontent.com/LunarVim/LunarVim/release-1.4/neovim-0.9/utils/installer/install.sh)
    - user: saltuser # Assuming 'saltuser' is the user for whom lvim is being installed
    - require:
      - pkg: lvim_install_prerequisites
