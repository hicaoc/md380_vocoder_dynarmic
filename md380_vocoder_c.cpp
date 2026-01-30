#include "md380_vocoder_c.h"
#include "md380_vocoder.h"
#include <string>
#include <cstring>

// Internal structure for vocoder instance
struct MD380VocoderInstance {
    MD380Vocoder* vocoder;
    std::string last_error;
    bool initialized;

    MD380VocoderInstance() : vocoder(nullptr), initialized(false) {
        vocoder = new MD380Vocoder();
    }

    ~MD380VocoderInstance() {
        if (vocoder) {
            delete vocoder;
            vocoder = nullptr;
        }
    }
};

extern "C" {

md380_vocoder_handle_t md380_vocoder_create(void) {
    try {
        MD380VocoderInstance* instance = new MD380VocoderInstance();
        return static_cast<md380_vocoder_handle_t>(instance);
    } catch (const std::exception& e) {
        return nullptr;
    }
}

void md380_vocoder_destroy(md380_vocoder_handle_t handle) {
    if (handle) {
        MD380VocoderInstance* instance = static_cast<MD380VocoderInstance*>(handle);
        delete instance;
    }
}

int md380_vocoder_init(md380_vocoder_handle_t handle) {
    if (!handle) {
        return -1;
    }

    MD380VocoderInstance* instance = static_cast<MD380VocoderInstance*>(handle);

    try {
        instance->vocoder->init();
        instance->initialized = true;
        instance->last_error.clear();
        return 0;
    } catch (const std::exception& e) {
        instance->last_error = std::string("Initialization failed: ") + e.what();
        return -1;
    } catch (...) {
        instance->last_error = "Initialization failed: unknown error";
        return -1;
    }
}

int md380_vocoder_encode(md380_vocoder_handle_t handle, const int16_t* pcm_in, uint8_t* ambe_out) {
    if (!handle || !pcm_in || !ambe_out) {
        return -1;
    }

    MD380VocoderInstance* instance = static_cast<MD380VocoderInstance*>(handle);

    if (!instance->initialized) {
        instance->last_error = "Vocoder not initialized";
        return -1;
    }

    try {
        instance->vocoder->encode(pcm_in, ambe_out);
        instance->last_error.clear();
        return 0;
    } catch (const std::exception& e) {
        instance->last_error = std::string("Encoding failed: ") + e.what();
        return -1;
    } catch (...) {
        instance->last_error = "Encoding failed: unknown error";
        return -1;
    }
}

int md380_vocoder_decode(md380_vocoder_handle_t handle, const uint8_t* ambe_in, int16_t* pcm_out) {
    if (!handle || !ambe_in || !pcm_out) {
        return -1;
    }

    MD380VocoderInstance* instance = static_cast<MD380VocoderInstance*>(handle);

    if (!instance->initialized) {
        instance->last_error = "Vocoder not initialized";
        return -1;
    }

    try {
        instance->vocoder->decode(ambe_in, pcm_out);
        instance->last_error.clear();
        return 0;
    } catch (const std::exception& e) {
        instance->last_error = std::string("Decoding failed: ") + e.what();
        return -1;
    } catch (...) {
        instance->last_error = "Decoding failed: unknown error";
        return -1;
    }
}

const char* md380_vocoder_get_error(md380_vocoder_handle_t handle) {
    if (!handle) {
        return "Invalid handle";
    }

    MD380VocoderInstance* instance = static_cast<MD380VocoderInstance*>(handle);

    if (instance->last_error.empty()) {
        return nullptr;
    }

    return instance->last_error.c_str();
}

} // extern "C"
