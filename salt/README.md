# SaltStack Project

This directory contains the SaltStack project for managing infrastructure.

## Structure

- `top.sls`: The main top file for state application.
- `pillar/top.sls`: The main top file for pillar data.
- `roles/`: Contains role definitions.
- `profile/`: Contains profile definitions.
- `formula/`: Contains Salt formulas (states).
  - `common/`: Common states applicable to all systems.
  - `_states/`: Custom state modules.
  - `_modules/`: Custom execution modules.
