
# FLEET_PROTOCOL System Configuration

This directory contains the configuration file for FLEET_PROTOCOL [CMCONF] system. It is used to set
required variables for [Package Tracker]. The most important variable is
`BA_PACKAGE_URI_TEMPLATE_REMOTE`, which sets URI of Package Repository used during builds.

This system must be installed in all Docker images.

More information about this system can be found in [Package Tracker Example].

Using this system is an example and in general Context does not need to use it.


[Package Tracker]: https://github.com/bacpack-system/package-tracker
[Package Tracker Example]: https://github.com/bacpack-system/package-tracker/tree/master/example
[CMCONF]: https://github.com/cmakelib/cmakelib-component-cmconf
