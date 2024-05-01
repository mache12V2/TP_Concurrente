#define N 100
#define WORKERS 4

mtype = {START, FINISH, CALCULATE}

chan request[WORKERS] = [0] of {byte, byte, byte}
chan response[WORKERS] = [0] of {int, int, int, int}

active proctype Worker(chan req, chan resp) {
    byte startIdx, endIdx
    int sumX, sumY, sumXY, sumXX

    do
    :: 
        req ? START, startIdx, endIdx ->
        sumX = 0
        sumY = 0
        sumXY = 0
        sumXX = 0
        for (i : startIdx+1 .. endIdx) {
            sumX += x[i]
            sumY += y[i]
            sumXY += x[i] * y[i]
            sumXX += x[i] * x[i]
        }
        resp ! sumX, sumY, sumXY, sumXX
    od
}

init {
    int i, j, k
    chan ch = [WORKERS] of {int, int, int, int}
    chan done = [WORKERS] of {byte}

    int x[N], y[N]
    for (i : 0 .. N-1) {
        x[i] = int(i)
        y[i] = 2 * int(i) + 5 + uniform(0, 10)
    }

    // Calculate
    for (i : 0 .. WORKERS-1) {
        run Worker(request[i], response[i])
    }

    for (i : 0 .. WORKERS-1) {
        startIdx = (N / WORKERS) * i
        endIdx = (N / WORKERS) * (i+1)
        if
        :: i == WORKERS-1 -> endIdx = N-1
        :: else -> skip
        fi
        request[i] ! START, startIdx, endIdx
    }

    for (i : 0 .. WORKERS-1) {
        response[i] ? sumX, sumY, sumXY, sumXX
        total.sumX += sumX
        total.sumY += sumY
        total.sumXY += sumXY
        total.sumXX += sumXX
    }

    int m, b
    m = (N * total.sumXY - total.sumX * total.sumY) / (N * total.sumXX - total.sumX * total.sumX)
    b = (total.sumY - m * total.sumX) / N
    printf("m = %g, b = %g\n", m, b)
}
