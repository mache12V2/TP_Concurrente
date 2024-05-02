#define N 1000
#define WORKERS 4

mtype = { CALCULATE, FINISH };

chan workerCh[WORKERS] = [0] of {mtype, int, int, int, int};
chan sem = [WORKERS] of {int};  // Canal semáforo para controlar la ejecución de los trabajadores

proctype worker(chan out; int start; int end) {
    sem ? 1; // Adquirir un "permiso" del semáforo

    int sumX = 0;
    int sumY = 0;
    int sumXY = 0;
    int sumXX = 0;
    int x;
    int y;
    int i = start;

    do
    :: i <= end -> 
        x = i; // Simulación de valores
        y = 2 * i + 5; // Simulación de valores
        sumX = sumX + x;
        sumY = sumY + y;
        sumXY = sumXY + x * y;
        sumXX = sumXX + x * x;
        i = i + 1;
    :: i > end -> break;
    od;

    out ! CALCULATE, sumX, sumY, sumXY, sumXX;

    sem ! 1; // Liberar el "permiso" del semáforo
}

init {
    int startIdx;
    int endIdx;
    int size = N / WORKERS;
    int i;

    // Inicializar el semáforo y lanzar trabajadores
    i = 0;
    do
    :: i < WORKERS - 1 -> 
        sem ! 1;  // Poner los "permisos" en el semáforo
        startIdx = i * size;
        endIdx = startIdx + size - 1;
        run worker(workerCh[i], startIdx, endIdx);
        i = i + 1;
    :: i == WORKERS - 1 ->  // Manejar el último trabajador aquí
        sem ! 1;
        startIdx = i * size;
        endIdx = N - 1;
        run worker(workerCh[i], startIdx, endIdx);
        i = i + 1;
    od;

    // Recolectar resultados
    int totalSumX = 0;
    int totalSumY = 0;
    int totalSumXY = 0;
    int totalSumXX = 0;
    mtype cmd;
    int px;
    int py;
    int pxy;
    int pxx;
    i = 0;

    do
    :: i < WORKERS -> 
        workerCh[i] ? cmd, px, py, pxy, pxx;
        totalSumX = totalSumX + px;
        totalSumY = totalSumY + py;
        totalSumXY = totalSumXY + pxy;
        totalSumXX = totalSumXX + pxx;
        i = i + 1;
    od;
    
    int m = (totalSumXY - totalSumX * totalSumY) / (totalSumXX - totalSumX * totalSumX);
    int b = (totalSumY - m * totalSumX);

    printf("m = %%d, b = %%d\\n", m, b);
}
