import pandas as pd
import matplotlib.pyplot as plt
import pandas as pd
import matplotlib.pyplot as plt
import sys
import time

def read_csv_with_retry(file_path, max_attempts=5, delay=1):
    for attempt in range(max_attempts):
        try:
            return pd.read_csv(file_path)
        except FileNotFoundError:
            if attempt < max_attempts - 1:
                time.sleep(delay)
            else:
                raise

try:
    df = read_csv_with_retry('points_Empleados.csv')
    m_values = df['m']
    b_values = df['b']
    
    plt.figure(figsize=(20, 12))
    plt.plot(m_values, b_values, marker='o', linestyle='-', color='b', label='Line Plot')
    
    plt.title('Recta de puntos m,b de los cálculos')
    plt.xlabel('m')
    plt.ylabel('b')
    plt.grid(True)
    plt.legend()
    plt.show()
    
    # Guarda el gráfico en un archivo temporal y muestra su nombre
    temp_file = 'temp_plot.png'
    plt.savefig(temp_file)
    
    # Imprime el nombre del archivo para que Go lo lea como respuesta
    print(temp_file)
except Exception as e:
    print(f"Error: {str(e)}", file=sys.stderr)
    sys.exit(1)
    
