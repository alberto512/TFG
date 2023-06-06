'use client';
import { Area, AreaChart, Tooltip, XAxis, YAxis } from 'recharts';
import styles from './AccountBalances.module.css';

type CustomTooltipProps = {
  active?: boolean;
  payload?: any[];
  label?: string;
};

const CustomTooltip = ({ active, payload, label }: CustomTooltipProps) => {
  if (active && payload && payload.length) {
    return (
      <div>
        <p>{label}</p>
        {payload.map((entry, index) => (
          <p key={index}>
            <span>{entry.name.charAt(0).toUpperCase() + entry.name.slice(1)}: </span>
            <span>{entry.value.toLocaleString('es')}</span>
          </p>
        ))}
      </div>
    );
  }

  return null;
};

type Props = {
  balances: any[];
  categories: string[];
};

const CategoriesComparison = ({ balances, categories }: Props) => {
  return (
    <div className={styles.categoriesComparison}>
      <AreaChart width={1000} height={250} data={balances}>
        <defs>
          <linearGradient id='colorCategory1' x1='0' y1='0' x2='0' y2='1'>
            <stop offset='5%' stopColor='#ff4136' stopOpacity={0.8} />
            <stop offset='95%' stopColor='#ff4136' stopOpacity={0} />
          </linearGradient>
          <linearGradient id='colorCategory2' x1='0' y1='0' x2='0' y2='1'>
            <stop offset='5%' stopColor='#82ca9d' stopOpacity={0.8} />
            <stop offset='95%' stopColor='#82ca9d' stopOpacity={0} />
          </linearGradient>
        </defs>
        <XAxis dataKey='date' tick={{ stroke: '#6b5b95' }} stroke={'#6b5b95'} />
        <YAxis
          padding={{ bottom: 15 }}
          tick={{ stroke: '#6b5b95' }}
          stroke={'#6b5b95'}
          tickFormatter={(value) => value.toLocaleString('es')}
        />
        <Tooltip content={(props) => <CustomTooltip {...props} />} />
        <Area type='monotone' dataKey={categories[0]} stroke='#ff4136' fillOpacity={1} fill='url(#colorCategory1)' />
        <Area type='monotone' dataKey={categories[1]} stroke='#82ca9d' fillOpacity={1} fill='url(#colorCategory2)' />
      </AreaChart>
    </div>
  );
};

export default CategoriesComparison;
