'use client';
import { Cell, Legend, Pie, PieChart } from 'recharts';
import styles from './CategoriesBalances.module.css';

type Balances = {
  amount: number;
  category: string;
};

type Props = {
  balancesPositive: Balances[];
  balancesNegative: Balances[];
};

const CategoriesBalances = ({ balancesPositive, balancesNegative }: Props) => {
  const COLORS_NEGATIVE = [
    '#FFA500',
    '#008000',
    '#9932CC',
    '#00FF00',
    '#EE82EE',
    '#FF0000',
    '#00CED1',
    '#FF4500',
    '#4169E1',
    '#00BFFF',
  ];

  const COLORS_POSITIVE = [
    '#0000FF',
    '#FF1493',
    '#FF4500',
    '#FFA07A',
    '#00BFFF',
    '#FF4500',
    '#9932CC',
    '#8B008B',
    '#32CD32',
    '#FFD700',
  ];

  return (
    <div className={styles.categoriesBalances}>
      {balancesNegative.length !== 0 && (
        <div className={styles.chartWrapper}>
          <span className={styles.label}>Negative balances</span>
          <PieChart width={400} height={400}>
            <Pie
              data={balancesNegative}
              dataKey='amount'
              nameKey='category'
              innerRadius={50}
              outerRadius={80}
              label={({ value }) => value.toLocaleString('de-DE')}
            >
              {balancesNegative.map((_, index) => (
                <Cell key={`${index}`} fill={COLORS_NEGATIVE[index % 10]} />
              ))}
            </Pie>
            <Legend />
          </PieChart>
        </div>
      )}
      {balancesPositive.length !== 0 && (
        <div className={styles.chartWrapper}>
          <span className={styles.label}>Postive balances</span>
          <PieChart width={400} height={400}>
            <Pie
              data={balancesPositive}
              dataKey='amount'
              nameKey='category'
              innerRadius={50}
              outerRadius={80}
              label={({ value }) => value.toLocaleString('de-DE')}
            >
              {balancesPositive.map((_, index) => (
                <Cell key={`${index}`} fill={COLORS_POSITIVE[index % 10]} />
              ))}
            </Pie>
            <Legend />
          </PieChart>
        </div>
      )}
    </div>
  );
};

export default CategoriesBalances;
