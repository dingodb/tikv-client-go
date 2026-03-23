# Copyright (c) 2024 dingodb.com, Inc. All Rights Reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.16)
project(tikvgo LANGUAGES C)

find_program(GO_EXECUTABLE go REQUIRED)

set(TIKVGO_LIB    ${CMAKE_CURRENT_BINARY_DIR}/libtikvgo.a)
set(TIKVGO_HEADER ${CMAKE_CURRENT_BINARY_DIR}/libtikvgo.h)

add_custom_command(
  OUTPUT ${TIKVGO_LIB} ${TIKVGO_HEADER}
  COMMAND ${CMAKE_COMMAND} -E env CGO_ENABLED=1
          ${GO_EXECUTABLE} build -buildmode=c-archive
          -o ${TIKVGO_LIB} .
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/bridge.go
          ${CMAKE_CURRENT_SOURCE_DIR}/go.mod
  COMMENT "Building TiKV Go CGO static library (libtikvgo.a)"
  VERBATIM
)

add_custom_target(tikvgo_build ALL DEPENDS ${TIKVGO_LIB} ${TIKVGO_HEADER})

add_library(tikvgo STATIC IMPORTED GLOBAL)
set_target_properties(tikvgo PROPERTIES
  IMPORTED_LOCATION ${TIKVGO_LIB}
)
add_dependencies(tikvgo tikvgo_build)

# Expose include directory for consumers when used as a subdirectory.
set(TIKVGO_INCLUDE_DIR ${CMAKE_CURRENT_BINARY_DIR} PARENT_SCOPE)

# Install rules: library and header.
install(FILES ${TIKVGO_LIB}    DESTINATION lib)
install(FILES ${TIKVGO_HEADER} DESTINATION include)
