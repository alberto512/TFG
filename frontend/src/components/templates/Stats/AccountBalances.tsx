'use client';
import { useState } from 'react';
import { Layer, Rectangle, Sankey } from 'recharts';
import InputSelect from '@elements/InputSelect/InputSelect';
import styles from './AccountBalances.module.css';

type Balances = {
  nodes: {
    name: string;
  }[];
  links: {
    source: number;
    target: number;
    value: number;
  }[];
};

type CustomNodeProps = {
  x: number;
  y: number;
  width: number;
  height: number;
  index: number;
  payload: any;
};

type Props = {
  balancesPositive: Balances;
  balancesNegative: Balances;
};

const CustomNode = ({ x, y, width, height, index, payload }: CustomNodeProps) => {
  const isOut = x + width > 500;

  const shouldShowLabels = height > 20;
  const fontSize = Math.min(Math.max(height / 2, 10), 30);

  const labelX = isOut ? x - 5 : x + width + 15;
  const labelY = y + height / 2 + fontSize / 2;

  return (
    <Layer key={`CustomNode${index}`}>
      <Rectangle x={x} y={y} width={width} height={height} fill='#8884d8' />
      {shouldShowLabels && (
        <>
          <text textAnchor={isOut ? 'end' : 'start'} x={labelX} y={labelY} fontSize={fontSize} fill='#000'>
            {isOut ? `${payload.name}: ${payload.value.toLocaleString('de-DE')}` : payload.value.toLocaleString('de-DE')}
          </text>
        </>
      )}
    </Layer>
  );
};

const AccountBalances = ({ balancesPositive, balancesNegative }: Props) => {
  const [statsTypeSelected, setStatsTypeSelected] = useState('0');
  const statTypes = [
    { id: '1', name: 'Negative Balances' },
    { id: '2', name: 'Positive Balances' },
  ];

  return (
    <div className={styles.accountBalances}>
      <div className={styles.selectWrapper}>
        <InputSelect placeholder='Select a type' options={statTypes} onChange={setStatsTypeSelected} />
      </div>
      {balancesNegative && statsTypeSelected === '1' && (
        <div className={styles.chartWrapper}>
          <Sankey
            width={1000}
            height={500}
            data={balancesNegative}
            nodePadding={5}
            node={(props) => <CustomNode {...props} />}
            link={{ stroke: '#dd002b' }}
          ></Sankey>
        </div>
      )}
      {balancesPositive && statsTypeSelected === '2' && (
        <div className={styles.chartWrapper}>
          <Sankey
            width={1000}
            height={500}
            data={balancesPositive}
            nodePadding={5}
            node={(props) => <CustomNode {...props} />}
            link={{ stroke: '#00ddb2' }}
          ></Sankey>
        </div>
      )}
    </div>
  );
};

export default AccountBalances;
