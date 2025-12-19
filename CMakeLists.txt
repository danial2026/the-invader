cmake_minimum_required(VERSION 3.10)
project(space_invaders_go)

# Custom target to call Go build
add_custom_target(game ALL
    COMMAND go build -o game ./cmd/invaders
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Building Go binary 'game'"
)
