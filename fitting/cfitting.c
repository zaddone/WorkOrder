#include "cfitting.h"
#include "math.h"

bool GetExceptVal(const double *ArrX,const int size,const double * Weights, const int cols){
	int endI = size-1;
	double a,x,dif, e,ex,dife;
	x = ArrX[endI];
	a = x;	
	ex = (ArrX[endI]+ (ArrX[endI]-ArrX[0]));
	e = ex;
	dif = Weights[0]+ a*Weights[1];
	dife = Weights[0]+ e*Weights[1];
	int i;      
	for (i = 2;i<cols;i++){
	        a*=x;
	        dif += a*Weights[i];

	        e*=ex;
	        dife += e*Weights[i];
	}
	return dife>dif;
//	return (ArrX[endI]-ArrX[0]) > 0;

}
bool CurveData(CvMat* Mx,CvMat* My,CvMat* Mw){

	if(Mx==NULL || My == NULL || Mw==NULL ) return false;
	CvMat *Mxs = cvCreateMat(Mx->cols,Mx->rows,CV_64FC1);
	cvTranspose(Mx,Mxs);  
	CvMat *Mul = cvCreateMat(Mxs->rows,Mx->cols,CV_64FC1);  
	cvMatMul(Mxs,Mx,Mul);
	CvMat *Mxx = cvCreateMat(Mul->rows,Mul->rows,CV_64FC1);
//	cvInvert(Mul,Mxx); 
	cvInvert(Mul,Mxx,CV_SVD);
	CvMat *Mxy = cvCreateMat(Mxs->rows,My->cols,CV_64FC1);
	cvMatMul(Mxs,My,Mxy);  
	
	cvMatMul(Mxx,Mxy,Mw);
	
	cvReleaseMat(&Mxs);
	cvReleaseMat(&Mul);
	cvReleaseMat(&Mxx);
	cvReleaseMat(&Mxy);
	return true;

}
unsigned int GetCurveWeight(const double *ArrX,const double *Y,const int len,const int xCols,double * Weights){
	CvMat* My = cvCreateMat(len,1,CV_64FC1);
	CvMat* Mx = cvCreateMat(len,xCols,CV_64FC1);
	int i,j,be;
	double a,_a;
	for(i = 0;i<len;i++){
		My->data.db[i] = Y[i];
		a = _a = ArrX[i];
		be = i * xCols;
		Mx->data.db[be] = 1;
		Mx->data.db[be+1] = a;
		for (j = 2; j < xCols ; j++){
			a = a * _a;
			if (isnan(a)){
				cvReleaseMat(&Mx);
				cvReleaseMat(&My);
				return 0;
			}
			Mx->data.db[be+j] = a;
		}
	}
	CvMat *Mw = cvCreateMat( xCols,1,CV_64FC1);
	CurveData(Mx,My,Mw);
	memcpy( Weights,Mw->data.db, sizeof(double)*xCols);
	cvReleaseMat(&Mx);
	cvReleaseMat(&My);
	cvReleaseMat(&Mw);
	return 1;
	
}

unsigned int GetCurveArr(const double *ArrX,const double *Y,const int len,const int Wsize){
	
	int be,Ws,i,n,j;
   	CvMat* My = cvCreateMat(len,1,CV_64FC1);
//	My->data.db = Y
	CvMat **Mxs =(CvMat **)malloc(sizeof(CvMat *)*Wsize);
	for(j = 0;j<Wsize;j++){
		Mxs[j] = cvCreateMat(len,j+2,CV_64FC1);
	}

	Ws = Wsize+2;
	double *X = (double *)malloc(sizeof(double)*Ws);
	X[0] = 1;
	double a;
	for(i = 0;i<len;i++){
		My->data.db[i] = Y[i];
		a = ArrX[i];
		X[1] = a;
		for (n = 2;n<Ws;n++){
			a*=X[1];
			X[n] = a;
		}
		for(j = 0;j<Wsize;j++){
		        be=i*Mxs[j]->cols;
		        for (n = 0;n<Mxs[j]->cols;n++){
		      		Mxs[j]->data.db[be+n] = X[n];
		        }
		}
	}
	free(X);
//	delete [] X; 
	unsigned int key = 0;
	int cols;
	for(j = 0;j<Wsize;j++){
		cols = Mxs[j]->cols;
		CvMat *Mw = cvCreateMat( cols,1,CV_64FC1);
		CurveData(Mxs[j],My,Mw);
		key = key<<1;
		if (GetExceptVal(ArrX,len,Mw->data.db,cols)){
			key++;
		}
		 
		cvReleaseMat(&Mxs[j]);
		cvReleaseMat(&Mw);
		
	}
	cvReleaseMat(&My);
//	delete [] ArrX;
//	free(ArrX);
	free(Mxs);
//	delete Mxs;
	return key;
}

unsigned int GetCurveData(const char *d,const int len,const int Wsize){
	int dl = sizeof(double);
	int l = len/2;
	int tl = l/dl;
//	printf("%d %d %d",len,l,tl);

	double *ArrX = (double*)malloc(l);
	memcpy(ArrX,d,l);

   	CvMat* My = cvCreateMat(tl,1,CV_64FC1);
	memcpy((My->data.db),d+l,l);

//	CvMat **Mxs = new CvMat *[Wsize];
	CvMat **Mxs =(CvMat **)malloc(sizeof(CvMat *)*Wsize);
	int j,i,n;
	for(j = 0;j<Wsize;j++){
		Mxs[j] = cvCreateMat(tl,j+2,CV_64FC1);
	}
	int Ws = Wsize+2;
	double *X = (double *)malloc(sizeof(double)*Ws);
	X[0] = 1;
	double a;
	int be;
	for(i = 0;i<tl;i++){
		a = ArrX[i];
		X[1] = a;
		for (n = 2;n<Ws;n++){
			a*=X[1];
			X[n] = a;
		}
		for(j = 0;j<Wsize;j++){
		        be=i*Mxs[j]->cols;
		        for (n = 0;n<Mxs[j]->cols;n++){
		      		Mxs[j]->data.db[be+n] = X[n];
		        }
		}
	}
	free(X);
//	delete [] X; 
	unsigned int key = 0;
	int cols;
	for(j = 0;j<Wsize;j++){
		cols = Mxs[j]->cols;
		CvMat *Mw = cvCreateMat( cols,1,CV_64FC1);
		CurveData(Mxs[j],My,Mw);
		key = key<<1;
		if (GetExceptVal(ArrX,tl,Mw->data.db,cols)){
			key++;
		}
		 
		cvReleaseMat(&Mxs[j]);
		cvReleaseMat(&Mw);
		
	}
	cvReleaseMat(&My);
//	delete [] ArrX;
	free(ArrX);
	free(Mxs);
//	delete Mxs;
	return key;

// 	CvMat* Mw = cvCreateMat(4,1,CV_64FC1);
//	printf("%d",len);
}
