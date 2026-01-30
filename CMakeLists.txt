cmake_minimum_required(VERSION 3.12)
include(FetchContent)
project(md380_vocoder)

set(CMAKE_CXX_STANDARD 20)

# Check if unzip and python3 are installed
find_program(UNZIP_COMMAND unzip)
if(NOT UNZIP_COMMAND)
    message(FATAL_ERROR "unzip not found")
endif()

find_program(PYTHON3_COMMAND python3)
if(NOT PYTHON3_COMMAND)
    message(FATAL_ERROR "python3 not found")
endif()

# Run download.sh in firmware directory, exit if it fails
execute_process(
    COMMAND sh download.sh
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/firmware)
if(NOT EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/firmware/D002.032.bin)
    message(FATAL_ERROR "download.sh failed")
endif()

add_custom_command(
    OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/firmware/firmware.c
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/firmware
    COMMAND xxd -i D002.032.bin | sed "s/unsigned char D002_032_bin/unsigned char firmware/g" | sed "s/unsigned int D002_032_bin_len/unsigned int firmware_len/g" > firmware.c
    COMMENT "Compiling firmware:D002.032.bin to binary")

add_custom_command(
    OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/firmware/sram.c
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/firmware
    COMMAND xxd -i d02032-core.img | sed "s/unsigned char d02032_core_img/unsigned char sram/g" | sed "s/unsigned int d02032_core_img_len/unsigned int sram_len/g" > sram.c
    COMMENT "Compiling sram:d02032-core.img to binary")

add_custom_target(firmware DEPENDS firmware/firmware.c firmware/sram.c)

# FetchContent_Declare(
#     dynarmic
#     GIT_REPOSITORY https://github.com/yuzu-mirror/dynarmic.git
#     GIT_TAG master
# )
# FetchContent_MakeAvailable(dynarmic)
add_subdirectory(dynarmic_local)

set_source_files_properties(${CMAKE_CURRENT_SOURCE_DIR}/firmware/firmware.c
                            PROPERTIES GENERATED TRUE)
set_source_files_properties(${CMAKE_CURRENT_SOURCE_DIR}/firmware/sram.c
                            PROPERTIES GENERATED TRUE)

# Build the main C++ vocoder library
add_library(md380_vocoder_cpp STATIC md380_vocoder.cpp ${CMAKE_CURRENT_SOURCE_DIR}/firmware/firmware.c ${CMAKE_CURRENT_SOURCE_DIR}/firmware/sram.c)
target_link_libraries(md380_vocoder_cpp STATIC dynarmic)

# Build the C wrapper library
add_library(md380_vocoder_c STATIC md380_vocoder_c.cpp)
target_link_libraries(md380_vocoder_c md380_vocoder_cpp)

# Create a combined static library with all dependencies for CGO
add_custom_command(
    OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libmd380_vocoder_full.a
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/c
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/cpp
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/dynarmic
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/fmt
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/mcl
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/zydis
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/zycore
    
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/c && ar x $<TARGET_FILE:md380_vocoder_c>
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/cpp && ar x $<TARGET_FILE:md380_vocoder_cpp>
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/dynarmic && ar x $<TARGET_FILE:dynarmic>
    
    # Using relative paths for dependencies
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/fmt && ar x ../../dynarmic_local/externals/fmt/libfmt.a
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/mcl && ar x ../../dynarmic_local/externals/mcl/src/libmcl.a
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/zydis && ar x ../../dynarmic_local/externals/zydis/libZydis.a
    COMMAND cd ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/zycore && ar x ../../dynarmic_local/externals/zydis/zycore/libZycore.a
    
    # Removing 'merge_temp/*.o' and adding explicit stored paths. 
    # Use 'find' or wildcard expansion in shell if possible, but cmake EXECUTE_PROCESS shell usage is tricky.
    # We can pass multiple wildcard paths to ar.
    COMMAND ar rcs ${CMAKE_CURRENT_BINARY_DIR}/libmd380_vocoder_full.a 
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/c/*.o
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/cpp/*.o
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/dynarmic/*.o
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/fmt/*.o
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/mcl/*.o
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/zydis/*.o
            ${CMAKE_CURRENT_BINARY_DIR}/merge_temp/zycore/*.o
            
    DEPENDS md380_vocoder_c md380_vocoder_cpp dynarmic
    COMMENT "Creating combined static library for CGO with isolated extractions"
)

add_custom_target(md380_vocoder_full ALL DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/libmd380_vocoder_full.a)

# Install headers and library to go_wrapper directory for CGO
install(FILES md380_vocoder_c.h DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}/go_wrapper/include)
install(FILES ${CMAKE_CURRENT_BINARY_DIR}/libmd380_vocoder_full.a
        DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}/go_wrapper/lib
        RENAME libmd380_vocoder.a)

# For backward compatibility, also create md380_vocoder target
add_library(md380_vocoder ALIAS md380_vocoder_cpp)

