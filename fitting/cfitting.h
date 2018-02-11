#pragma once
#include <stdio.h>
#include <stdbool.h>
#include <opencv/cv.h>
#include <opencv/highgui.h> 

bool CurveData(CvMat* Mx,CvMat* My,CvMat* Mw);
unsigned int GetCurveWeight(const double *ArrX,const double *Y,const int len,const int xCols,double * Weights);

bool GetExceptVal(const double *ArrX,const int size,const double * Weights, const int cols);
unsigned int GetCurveData(const char *d,const int len,const int Wsize);
unsigned int GetCurveArr(const double *ArrX,const double *Y,const int len,const int Wsize);
