#ifndef MD380_VOCODER_C_H
#define MD380_VOCODER_C_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// Opaque handle type for vocoder instance
typedef void* md380_vocoder_handle_t;

/**
 * Create a new vocoder instance
 * Returns NULL on failure
 */
md380_vocoder_handle_t md380_vocoder_create(void);

/**
 * Destroy a vocoder instance
 */
void md380_vocoder_destroy(md380_vocoder_handle_t handle);

/**
 * Initialize the vocoder
 * Returns 0 on success, non-zero on failure
 */
int md380_vocoder_init(md380_vocoder_handle_t handle);

/**
 * Encode PCM to AMBE
 * pcm_in: 160 int16_t samples
 * ambe_out: 9 bytes output buffer
 * Returns 0 on success, non-zero on failure
 */
int md380_vocoder_encode(md380_vocoder_handle_t handle, const int16_t* pcm_in, uint8_t* ambe_out);

/**
 * Decode AMBE to PCM
 * ambe_in: 9 bytes input
 * pcm_out: 160 int16_t samples output buffer
 * Returns 0 on success, non-zero on failure
 */
int md380_vocoder_decode(md380_vocoder_handle_t handle, const uint8_t* ambe_in, int16_t* pcm_out);

/**
 * Get last error message
 * Returns error string or NULL if no error
 */
const char* md380_vocoder_get_error(md380_vocoder_handle_t handle);

#ifdef __cplusplus
}
#endif

#endif // MD380_VOCODER_C_H
