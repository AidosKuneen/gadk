// +build amd64

package gadk

// #cgo LDFLAGS:-lpthread
// #cgo CPPFLAGS: -msse2 -Wall
/*
 #include <stdio.h>
 #include <stdlib.h>
 #include <string.h>
 #include <time.h>
 #include <stddef.h>

 #ifdef _MSC_VER
 #include <intrin.h>
 #else
 #include <x86intrin.h>
 #endif

 #if defined(__linux) || defined(__APPLE__)
 #include <unistd.h>
 #elif defined(__MINGW64__)
 #include <windows.h>
 #include <winbase.h>
 #elif _WIN32
 #include <windows.h>
 #include <process.h>
 #endif

 #if defined(_MSC_VER) || defined(__MINGW64__)
 #include <malloc.h>
 #endif // defined(_MSC_VER) || defined(__MINGW32__)

 #ifndef _MSC_VER
 #include <pthread.h>
 #endif

 #ifdef _WIN32
 HANDLE	mutex = NULL;
 #else
 pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
 #endif

#if !defined(__cplusplus)
#  ifdef _MSC_VER
#    define alignas(n)  __declspec(align(n))
#  else
#    define alignas(n)  __attribute__((aligned(n)))
#  endif  // _MSC_VER
#endif


 #define HBITS 0xFFFFFFFFFFFFFFFFuLL
 #define LBITS 0x0000000000000000uLL
 #define HASH_LENGTH 243 //trits
 #define STATE_LENGTH 3 * HASH_LENGTH //trits
 #define HALF_LENGTH 364 //trits
 #define TX_LENGTH 2673 //trytes

 #define LOW00 0xDB6DB6DB6DB6DB6DuLL //0b1101101101101101101101101101101101101101101101101101101101101101L;
 #define HIGH00 0xB6DB6DB6DB6DB6DBuLL //0b1011011011011011011011011011011011011011011011011011011011011011L;
 #define LOW10 0xF1F8FC7E3F1F8FC7uLL //0b1111000111111000111111000111111000111111000111111000111111000111L;
 #define HIGH10 0x8FC7E3F1F8FC7E3FuLL //0b1000111111000111111000111111000111111000111111000111111000111111L;
 #define LOW20 0x7FFFE00FFFFC01FFuLL //0b0111111111111111111000000000111111111111111111000000000111111111L;
 #define HIGH20 0xFFC01FFFF803FFFFuLL //0b1111111111000000000111111111111111111000000000111111111111111111L;
 #define LOW30 0xFFC0000007FFFFFFuLL //0b1111111111000000000000000000000000000111111111111111111111111111L;
 #define HIGH30 0x003FFFFFFFFFFFFFuLL //0b0000000000111111111111111111111111111111111111111111111111111111L;
 #define LOW40 0xFFFFFFFFFFFFFFFFuLL //0b1111111111111111111111111111111111111111111111111111111111111111L;
 #define HIGH40 0xFFFFFFFFFFFFFFFFuLL //0b1111111111111111111111111111111111111111111111111111111111111111L;

 #define LOW01 0x6DB6DB6DB6DB6DB6uLL //0b0110110110110110110110110110110110110110110110110110110110110110
 #define HIGH01 0xDB6DB6DB6DB6DB6DuLL //0b1101101101101101101101101101101101101101101101101101101101101101
 #define LOW11 0xF8FC7E3F1F8FC7E3uLL //0b1111100011111100011111100011111100011111100011111100011111100011
 #define HIGH11 0xC7E3F1F8FC7E3F1FuLL //0b1100011111100011111100011111100011111100011111100011111100011111
 #define LOW21 0xC01FFFF803FFFF00uLL //0b1100000000011111111111111111100000000011111111111111111100000000
 #define HIGH21 0x3FFFF007FFFE00FFuLL //0b0011111111111111111100000000011111111111111111100000000011111111
 #define LOW31 0x00000FFFFFFFFFFFuLL //0b0000000000000000000011111111111111111111111111111111111111111111
 #define HIGH31 0xFFFFFFFFFFFE0000uLL //0b1111111111111111111111111111111111111111111111100000000000000000
 #define LOW41 0x000000000001FFFFuLL //0b0000000000000000000000000000000000000000000000011111111111111111
 #define HIGH41 0xFFFFFFFFFFFFFFFFuLL //0b1111111111111111111111111111111111111111111111111111111111111111

 #ifndef EXPORT
 #if defined(_WIN32) || defined(_MSC_VER)
 #define EXPORT __declspec(dllexport)
 #else
 #define EXPORT
 #endif
 #endif

 typedef struct{
     int corenum;
     long long count;
     long long time;
     char * trytes;
 } RESULT;

 #ifdef __cplusplus
 extern "C"{
 #endif
     EXPORT void ccurl_pow_finalize();
     EXPORT void ccurl_pow_interrupt();
     EXPORT char* ccurl_pow(char* trytes, int minWeightMagnitude, RESULT *result);
 #ifdef __cplusplus
 }
 #endif

 int getCpuNum()
 {
 #if defined(__linux) || defined(__APPLE__)
     // for linux
     return sysconf(_SC_NPROCESSORS_ONLN);
 #elif defined(__MINGW64__) || defined(__MINGW32__) || defined(_MSC_VER)
     // for windows and wine
     SYSTEM_INFO info;
     GetSystemInfo(&info);
     return info.dwNumberOfProcessors;
 #endif
 }

 const int indices[] = {
     0, 364, 728, 363, 727, 362, 726, 361, 725, 360, 724, 359, 723, 358, 722, 357, 721, 356, 720, 355, 719, 354, 718, 353, 717, 352, 716, 351, 715, 350, 714, 349, 713, 348, 712, 347, 711, 346, 710, 345, 709, 344, 708, 343, 707, 342, 706, 341, 705, 340, 704, 339, 703, 338, 702, 337, 701, 336, 700, 335, 699, 334, 698, 333, 697, 332, 696, 331, 695, 330, 694, 329, 693, 328, 692, 327, 691, 326, 690, 325, 689, 324, 688, 323, 687, 322, 686, 321, 685, 320, 684, 319, 683, 318, 682, 317, 681, 316, 680, 315, 679, 314, 678, 313, 677, 312, 676, 311, 675, 310, 674, 309, 673, 308, 672, 307, 671, 306, 670, 305, 669, 304, 668, 303, 667, 302, 666, 301, 665, 300, 664, 299, 663, 298, 662, 297, 661, 296, 660, 295, 659, 294, 658, 293, 657, 292, 656, 291, 655, 290, 654, 289, 653, 288, 652, 287, 651, 286, 650, 285, 649, 284, 648, 283, 647, 282, 646, 281, 645, 280, 644, 279, 643, 278, 642, 277, 641, 276, 640, 275, 639, 274, 638, 273, 637, 272, 636, 271, 635, 270, 634, 269, 633, 268, 632, 267, 631, 266, 630, 265, 629, 264, 628, 263, 627, 262, 626, 261, 625, 260, 624, 259, 623, 258, 622, 257, 621, 256, 620, 255, 619, 254, 618, 253, 617, 252, 616, 251, 615, 250, 614, 249, 613, 248, 612, 247, 611, 246, 610, 245, 609, 244, 608, 243, 607, 242, 606, 241, 605, 240, 604, 239, 603, 238, 602, 237, 601, 236, 600, 235, 599, 234, 598, 233, 597, 232, 596, 231, 595, 230, 594, 229, 593, 228, 592, 227, 591, 226, 590, 225, 589, 224, 588, 223, 587, 222, 586, 221, 585, 220, 584, 219, 583, 218, 582, 217, 581, 216, 580, 215, 579, 214, 578, 213, 577, 212, 576, 211, 575, 210, 574, 209, 573, 208, 572, 207, 571, 206, 570, 205, 569, 204, 568, 203, 567, 202, 566, 201, 565, 200, 564, 199, 563, 198, 562, 197, 561, 196, 560, 195, 559, 194, 558, 193, 557, 192, 556, 191, 555, 190, 554, 189, 553, 188, 552, 187, 551, 186, 550, 185, 549, 184, 548, 183, 547, 182, 546, 181, 545, 180, 544, 179, 543, 178, 542, 177, 541, 176, 540, 175, 539, 174, 538, 173, 537, 172, 536, 171, 535, 170, 534, 169, 533, 168, 532, 167, 531, 166, 530, 165, 529, 164, 528, 163, 527, 162, 526, 161, 525, 160, 524, 159, 523, 158, 522, 157, 521, 156, 520, 155, 519, 154, 518, 153, 517, 152, 516, 151, 515, 150, 514, 149, 513, 148, 512, 147, 511, 146, 510, 145, 509, 144, 508, 143, 507, 142, 506, 141, 505, 140, 504, 139, 503, 138, 502, 137, 501, 136, 500, 135, 499, 134, 498, 133, 497, 132, 496, 131, 495, 130, 494, 129, 493, 128, 492, 127, 491, 126, 490, 125, 489, 124, 488, 123, 487, 122, 486, 121, 485, 120, 484, 119, 483, 118, 482, 117, 481, 116, 480, 115, 479, 114, 478, 113, 477, 112, 476, 111, 475, 110, 474, 109, 473, 108, 472, 107, 471, 106, 470, 105, 469, 104, 468, 103, 467, 102, 466, 101, 465, 100, 464, 99, 463, 98, 462, 97, 461, 96, 460, 95, 459, 94, 458, 93, 457, 92, 456, 91, 455, 90, 454, 89, 453, 88, 452, 87, 451, 86, 450, 85, 449, 84, 448, 83, 447, 82, 446, 81, 445, 80, 444, 79, 443, 78, 442, 77, 441, 76, 440, 75, 439, 74, 438, 73, 437, 72, 436, 71, 435, 70, 434, 69, 433, 68, 432, 67, 431, 66, 430, 65, 429, 64, 428, 63, 427, 62, 426, 61, 425, 60, 424, 59, 423, 58, 422, 57, 421, 56, 420, 55, 419, 54, 418, 53, 417, 52, 416, 51, 415, 50, 414, 49, 413, 48, 412, 47, 411, 46, 410, 45, 409, 44, 408, 43, 407, 42, 406, 41, 405, 40, 404, 39, 403, 38, 402, 37, 401, 36, 400, 35, 399, 34, 398, 33, 397, 32, 396, 31, 395, 30, 394, 29, 393, 28, 392, 27, 391, 26, 390, 25, 389, 24, 388, 23, 387, 22, 386, 21, 385, 20, 384, 19, 383, 18, 382, 17, 381, 16, 380, 15, 379, 14, 378, 13, 377, 12, 376, 11, 375, 10, 374, 9, 373, 8, 372, 7, 371, 6, 370, 5, 369, 4, 368, 3, 367, 2, 366, 1, 365, 0
 };

 const char truthTable128[] = { 1, 0, -1, 0, 1, -1, 0, 0, -1, 1, 0 };

 const char tryte2trit_table[][3] = {
     { 0, 0, 0 }, { 1, 0, 0 }, { -1, 1, 0 }, { 0, 1, 0 }, { 1, 1, 0 }, { -1, -1, 1 }, { 0, -1, 1 }, { 1, -1, 1 }, { -1, 0, 1 }, { 0, 0, 1 }, { 1, 0, 1 }, { -1, 1, 1 }, { 0, 1, 1 }, { 1, 1, 1 }, { -1, -1, -1 }, { 0, -1, -1 }, { 1, -1, -1 }, { -1, 0, -1 }, { 0, 0, -1 }, { 1, 0, -1 }, { -1, 1, -1 }, { 0, 1, -1 }, { 1, 1, -1 }, { -1, -1, 0 }, { 0, -1, 0 }, { 1, -1, 0 }, { -1, 0, 0 }
 };

 static int stop = 0;

 void tx2trit(char tryte[], char trit[])
 {
     int i = 0;
     if (strlen(tryte) != TX_LENGTH) {
         fprintf(stderr, "illegal tx length.\n");
         exit(1);
     }
     for (i = 0; i < TX_LENGTH; i++) {
         if (!(tryte[i] == '9' || (tryte[i] >= 'A' && tryte[i] <= 'Z'))) {
             fprintf(stderr, "illegal tryte.\n");
             exit(1);
         }
         int ch = 0;
         if (tryte[i] != '9') {
             ch = tryte[i] - 'A' + 1;
         }
         trit[i * 3] = tryte2trit_table[ch][0];
         trit[i * 3 + 1] = tryte2trit_table[ch][1];
         trit[i * 3 + 2] = tryte2trit_table[ch][2];
     }
 }

 void hash2tryte(char trit[], char tryte[])
 {
     int i = 0;
     for (i = 0; i < HASH_LENGTH / 3; i++) {
         int n = trit[i * 3] + trit[i * 3 + 1] * 3 + trit[i * 3 + 2] * 9;
         if (n < 0) {
             n += 27;
         }
         if (n > 27 || n < 0) {
             fprintf(stderr, "illegal trit. %d\n", n);
             exit(1);
         }
         if (n == 0) {
             tryte[i] = '9';
         }
         else {
             tryte[i] = (char)n + 'A' - 1;
         }
     }
     tryte[HASH_LENGTH / 3] = 0;
 }

 void transform128(char state[])
 {
     int r = 0, i = 0;
     char copy[STATE_LENGTH] = { 0 };
     for (r = 0; r < 27; r++) {
         memcpy(copy, state, STATE_LENGTH);
         for (i = 0; i < STATE_LENGTH; i++) {
             int aa = indices[i];
             int bb = indices[i + 1];
             state[i] = truthTable128[copy[aa] + (copy[bb] << 2) + 5];
         }
     }
 }

 void absorb(char state[], char in[], int size)
 {
     int len = 0, offset = 0;
     for (offset = 0; offset < size; offset += len) {
         len = 243;
         if (size < 243) {
             len = size;
         }
         memcpy(state, in + offset, len);
         transform128(state);
     }
 }

 void transform64_128(__m128i* lmid, __m128i* hmid)
 {
     int j, r, t1, t2;
       __m128i one = _mm_set_epi64x(HBITS, HBITS);
     __m128i alpha, beta, gamma, delta,nalpha;
     __m128i *lto = lmid + STATE_LENGTH, *hto = hmid + STATE_LENGTH;
     __m128i *lfrom = lmid, *hfrom = hmid;
     for (r = 0; r < 26; r++) {
         for (j = 0; j < STATE_LENGTH; j++) {
             t1 = indices[j];
             t2 = indices[j + 1];

             alpha = lfrom[t1];
             beta = hfrom[t1];
             gamma = hfrom[t2];
             //(alpha | (~gamma)) & (lfrom[t2] ^ beta);
             nalpha = _mm_andnot_si128(alpha,gamma);
             delta =  _mm_andnot_si128(nalpha, _mm_xor_si128(lfrom[t2],beta));

             lto[j] = _mm_andnot_si128(delta,one);  //~delta;
             hto[j] = _mm_or_si128(_mm_xor_si128(alpha,gamma),delta); //(alpha ^ gamma) | delta;
         }
         __m128i *lswap = lfrom, *hswap = hfrom;
         lfrom = lto;
         hfrom = hto;
         lto = lswap;
         hto = hswap;
     }
     for (j = 0; j < HASH_LENGTH; j++) {
         t1 = indices[j];
         t2 = indices[j + 1];

         alpha = lfrom[t1];
         beta = hfrom[t1];
         gamma = hfrom[t2];
         nalpha = _mm_andnot_si128(alpha,gamma);
         delta =  _mm_andnot_si128(nalpha, _mm_xor_si128(lfrom[t2],beta));
         lto[j] = _mm_andnot_si128(delta,one);  //~delta;
             hto[j] = _mm_or_si128(_mm_xor_si128(alpha,gamma),delta); //(alpha ^ gamma) | delta;
      }
 }

 int incr128(__m128i* mid_low, __m128i* mid_high)
 {
     int i;
     __m128i carry;
     alignas(16) unsigned long long c[2]={0};
     for (i = 5; i < HASH_LENGTH && (i == 5 || c[0] ); i++) {
         __m128i low = mid_low[i], high = mid_high[i];
         mid_low[i] = _mm_xor_si128(high,low);
         mid_high[i] = low;
         carry = _mm_andnot_si128(low,high);
        _mm_store_si128 ((__m128i*)c,carry);
     }
     return i == HASH_LENGTH;
 }

 void seri128(__m128i* low, __m128i* high, int n, char* r)
 {
     int i = 0, index = 0;
     alignas(16) unsigned long long c[2]={0};

     if (n > 63) {
         n -= 64;
         index = 1;
     }
     for (i = 0; i < HASH_LENGTH; i++) {
        _mm_store_si128 ((__m128i*)c,low[i]);
         unsigned long long ll = (c[index] >> n) & 1;
        _mm_store_si128 ((__m128i*)c,high[i]);
         unsigned long long hh = (c[index] >> n) & 1;
         if (hh == 0 && ll == 1) {
             r[i] = -1;
         }
         if (hh == 1 && ll == 1) {
             r[i] = 0;
         }
         if (hh == 1 && ll == 0) {
             r[i] = 1;
         }
     }
 }

 int readStop(){
     int stop_=0;
 #ifdef _WIN32
     WaitForSingleObject(mutex, INFINITE );
 #else
     pthread_mutex_lock(&mutex);
 #endif
     stop_=stop;
 #ifdef _WIN32
    ReleaseMutex(mutex);
 #else
     pthread_mutex_unlock(&mutex);
 #endif
     return stop_;
 }

 void writeStop(int stop_){
 #ifdef _WIN32
     WaitForSingleObject(mutex, INFINITE );
 #else
     pthread_mutex_lock(&mutex);
 #endif
     stop=stop_;
 #ifdef _WIN32
    ReleaseMutex(mutex);
 #else
     pthread_mutex_unlock(&mutex);
 #endif
 }

 int check128(__m128i* l, __m128i* h, int m)
 {
     int i; //omit init for speed
     alignas(16) unsigned long long c[2]={0};

     __m128i nonce_probe = _mm_set_epi64x(HBITS, HBITS);
     for (i = HASH_LENGTH - m; i < HASH_LENGTH; i++) {
          nonce_probe =  _mm_andnot_si128(_mm_xor_si128(l[i],h[i]),nonce_probe);
         _mm_store_si128 ((__m128i*)c,nonce_probe);
         if (c[0] == LBITS && c[1] == LBITS) {
             return -1;
         }
     }
      _mm_store_si128 ((__m128i*)c,nonce_probe);
     for (i = 0; i < 64; i++) {
         if ((c[0] >> i) & 1) {
             return i + 0 * 64;
         }
         if ((c[1] >> i) & 1) {
             return i + 1 * 64;
         }
     }
     return -2;
 }

 long long int loop_cpu128(__m128i* lmid, __m128i* hmid, int m, char* nonce)
 {
     int n = 0, j = 0;
     long long int i = 0;
     __m128i lcpy[STATE_LENGTH * 2], hcpy[STATE_LENGTH * 2];
     for (i = 0; !incr128(lmid, hmid); i++) {
         for (j = 0; j < STATE_LENGTH; j++) {
             lcpy[j] = lmid[j];
             hcpy[j] = hmid[j];
         }
         transform64_128(lcpy, hcpy);
         if ((n = check128(lcpy + STATE_LENGTH, hcpy + STATE_LENGTH, m)) >= 0) {
             seri128(lmid, hmid, n, nonce);
             return i * 128;
         }
         int stop_=readStop();
         if (stop_) {
             return -i * 128 - 1;
         }
     }
     return -i * 128 - 1;
 }

 // 01:-1 11:0 10:1
 void para128(char in[], __m128i l[], __m128i h[])
 {
     int i = 0;
     for (i = 0; i < STATE_LENGTH; i++) {
         switch (in[i]) {
         case 0:
             l[i] = _mm_set_epi64x(HBITS, HBITS);
             h[i] = _mm_set_epi64x(HBITS, HBITS);
             break;
         case 1:
             l[i] = _mm_set_epi64x(LBITS, LBITS);
             h[i] = _mm_set_epi64x(HBITS, HBITS);
             break;
         case -1:
             l[i] = _mm_set_epi64x(HBITS, HBITS);
             h[i] = _mm_set_epi64x(LBITS, LBITS);
             break;
         }
     }
 }



 void incrN128(int n, __m128i* mid_low, __m128i* mid_high)
 {
     int i, j;
     alignas(16) unsigned long long c[2]={1,1};
     __m128i carry;
     for (j = 0; j < n; j++) {
         carry = _mm_set_epi64x(HBITS, HBITS);
         for (i = HASH_LENGTH - 7; i < HASH_LENGTH && c[0] ; i++) {
             __m128i low = mid_low[i], high = mid_high[i];
             mid_low[i] = _mm_xor_si128(high,low);
             mid_high[i] = low;
             carry = _mm_andnot_si128(low,high);
             _mm_store_si128 ((__m128i*)c,carry);
         }
     }
 }

 typedef struct param {
     char* mid;
     int mwm;
     char nonce[HASH_LENGTH];
     int n;
     long long int count;
 } PARAM;

 #ifdef _WIN32
 unsigned __stdcall  pwork_(void* p)
 #else
 void *pwork_(void* p)
 #endif
 {
     PARAM* par = (PARAM*)(p);
     __m128i lmid[STATE_LENGTH], hmid[STATE_LENGTH];

     para128(par->mid, lmid, hmid);
     lmid[0] = _mm_set_epi64x(LOW00, LOW01);
     hmid[0] = _mm_set_epi64x(HIGH00, HIGH01);
     lmid[1] = _mm_set_epi64x(LOW10, LOW11);
     hmid[1] = _mm_set_epi64x(HIGH10, HIGH11);
     lmid[2] = _mm_set_epi64x(LOW20, LOW21);
     hmid[2] = _mm_set_epi64x(HIGH20, HIGH21);
     lmid[3] = _mm_set_epi64x(LOW30, LOW31);
     hmid[3] = _mm_set_epi64x(HIGH30, HIGH31);
     lmid[4] = _mm_set_epi64x(LOW40, LOW41);
     hmid[4] = _mm_set_epi64x(HIGH40, HIGH41);

     incrN128(par->n, lmid, hmid);
     par->count = loop_cpu128(lmid, hmid, par->mwm, par->nonce);
     if (par->count >= 0) {
         writeStop(1);
     }
 #ifdef _WIN32
     return 0;
 #else
     return NULL;
 #endif
 }

 long long int pwork128(char tx[], int mwm, char nonce[])
 {
 #ifdef _WIN32
	mutex = CreateMutex(NULL,FALSE,NULL);
 #endif
     long long int countSSE = 0;
     int i = 0;
     char trits[TX_LENGTH * 3] = { 0 }, mid[STATE_LENGTH] = { 0 };

     tx2trit(tx, trits);
     absorb(mid, trits, TX_LENGTH * 3 - HASH_LENGTH);
     int procs = getCpuNum();
     if (procs>1){
         // procs--;
     }
 #ifdef _WIN32
     HANDLE *thread = (HANDLE*)calloc(sizeof(HANDLE), procs);
 #else
     pthread_t* thread = (pthread_t*)calloc(sizeof(pthread_t), procs);
 #endif
     PARAM* p = (PARAM*)calloc(sizeof(PARAM), procs);
     for (i = 0; i < procs; i++) {
         p[i].mid = mid;
         p[i].mwm = mwm;
         p[i].n = i;
 #ifdef _WIN32
         thread[i] = (HANDLE)_beginthreadex(NULL, 0, pwork_, (LPVOID)&p[i], 0, NULL);
         if (thread[i]==NULL) {
 #else
         int ret = pthread_create(&thread[i], NULL, pwork_, &p[i]);
         if (ret != 0) {
 #endif
             fprintf(stderr, "can not create thread\n");
             return -1;
         }
     }
     for (i = 0; i < procs; i++) {
 #ifdef _WIN32
         int ret = WaitForSingleObject( thread[i], INFINITE );
         CloseHandle(thread[i]);
         if (ret == WAIT_FAILED) {
 #else
         int ret = pthread_join(thread[i], NULL);
         if (ret != 0) {
 #endif
             fprintf(stderr, "can not join thread\n");
             return -1;
         }
         if (p[i].count >= 0) {
             memcpy(nonce, p[i].nonce, HASH_LENGTH);
             countSSE += p[i].count;
         }
         else {
             countSSE += -p[i].count + 1;
         }
     }
     free(thread);
     free(p);

     return countSSE;
 }
*/
import "C"
import (
	"errors"
	"unsafe"
)

func init() {
	pows["PowSSE"] = PowSSE
}

var countSSE int64

//PowSSE is proof of work of adk for amd64 using SSE2(or AMD64).
func PowSSE(trytes Trytes, mwm int) (Trytes, error) {
	C.writeStop(0)
	defer func() {
		C.writeStop(0)
	}()
	nonce := make(Trits, HashSize)
	tr := []byte(trytes)
	cnt := C.pwork128((*C.char)(unsafe.Pointer(&tr[0])), C.int(mwm), (*C.char)(unsafe.Pointer(&nonce[0])))
	countSSE = int64(cnt)
	if cnt >= 0 {
		result := nonce.Trytes()
		return result, nil
	}
	return "", errors.New("failed to PoW")
}
