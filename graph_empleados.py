import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_csv('points_Empleados.csv')
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