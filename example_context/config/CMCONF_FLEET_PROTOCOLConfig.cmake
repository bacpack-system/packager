
# Example configuration for CMCONF system.
#

FIND_PACKAGE(CMLIB REQUIRED COMPONENTS CMCONF)

CMCONF_INIT_SYSTEM(FLEET_PROTOCOL)

#
# Setting using upstream Package Repository by default for this system. This can be overridden by
# App in CMakeLists.txt.
#
CMCONF_SET(BA_PACKAGE_LOCAL_USE OFF)
CMCONF_SET(BA_PACKAGE_LOCAL_PATH "")

#
# The http authorization header is usually used for accessing private Package Repositories. This
# example does not need it, but the variable must be set.
#
CMCONF_SET(BA_PACKAGE_HTTP_AUTHORIZATION_HEADER "")

#
# Setting BringAuto's Package Repository URI Template
#
CMCONF_SET(BA_PACKAGE_URI_REVISION master)
CMCONF_SET(BA_PACKAGE_URI_TEMPLATE_REMOTE "https://gitea.bringauto.com/fleet-protocol/package-repository/media/<REVISION>/package/<GIT_PATH>/<PACKAGE_GROUP_NAME>/<ARCHIVE_NAME>")

#
# Gitea hosted public Package Repository:
#
#CMCONF_SET(BA_PACKAGE_URI_TEMPLATE_REMOTE "https://gitea.example.com/username/repository/media/<REVISION>/package/<GIT_PATH>/<PACKAGE_GROUP_NAME>/<ARCHIVE_NAME>")

#
# Gitea hosted private Package Repository.
# Do not forget to specify Access Token
#
#CMCONF_SET(BA_PACKAGE_HTTP_AUTHORIZATION_HEADER "token <token>")
#CMCONF_SET(BA_PACKAGE_URI_TEMPLATE_REMOTE "https://gitea.example.com/username/repository/raw/<REVISION>/package/<GIT_PATH>/<PACKAGE_GROUP_NAME>/<ARCHIVE_NAME>")

#
# Gitlab hosted private Package Repository.
#
#CMCONF_SET(BA_PACKAGE_HTTP_AUTHORIZATION_HEADER "Bearer <token>")
#CMCONF_SET(BA_PACKAGE_URI_TEMPLATE_REMOTE "https://gitlab.example.com/username/repository/-/raw/<REVISION>/package/<GIT_PATH>/<PACKAGE_GROUP_NAME>/<ARCHIVE_NAME>")
