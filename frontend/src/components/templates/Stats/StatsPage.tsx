'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import CategoriesBalances from './CategoriesBalances';
import AccountBalances from './AccountBalances';
import CategoriesComparison from './CategoriesComparison';
import InputSelect from '@elements/InputSelect/InputSelect';
import styles from './StatsPage.module.css';
import Loading from '@modules/Loading/Loading';

type Balances1 = {
  amount: number;
  category: string;
};

type Balances2 = {
  nodes: {
    name: string;
  }[];
  links: {
    source: number;
    target: number;
    value: number;
  }[];
};

const StatsPage = () => {
  const [statType, setStatType] = useState('');
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [accountSelected, setAccountSelected] = useState<string>('');
  const [accountsSelected, setAccountsSelected] = useState<string[]>([]);
  const [categoriesSelected, setCategoriesSelected] = useState<string[]>([]);
  const [categoriesPositionSelected, setCategoriesPositionSelected] = useState({
    0: '',
    1: '',
  });
  const [statsReady, setStatsReady] = useState(false);
  const [loading, setLoading] = useState(false);
  const [balancesPositive1, setBalancesPositive1] = useState<Balances1[]>([]);
  const [balancesNegative1, setBalancesNegative1] = useState<Balances1[]>([]);
  const [balancesPositive2, setBalancesPositive2] = useState<Balances2>({ nodes: [], links: [] });
  const [balancesNegative2, setBalancesNegative2] = useState<Balances2>({ nodes: [], links: [] });
  const [balances3, setBalances3] = useState<any[]>([]);
  const [categoriesComparison, setCategoriesComparison] = useState<string[]>([]);

  const router = useRouter();

  const statTypes = [
    { id: '1', name: 'Categories Balances' },
    { id: '2', name: 'Account Balances' },
    { id: '3', name: 'Categories Comparison' },
  ];

  const handleStatTypeChanged = () => {
    setAccountsSelected([]);
    setCategoriesSelected([]);
    axios
      .get('/api/accounts', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        setAccounts(response.data.map((account: Account) => ({ id: account.id, name: account.iban })));
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });

    axios
      .get('/api/categories', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        setCategories(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  const handleStat1Changed = () => {
    setBalancesNegative1([]);
    setBalancesPositive1([]);
    axios
      .post(
        '/api/stats',
        {
          selectedAccounts: accountsSelected,
          selectedCategories: categoriesSelected,
          initDate: new Date('2020-01-01').getTime(),
          endDate: new Date('2020-12-31').getTime(),
        },
        {
          headers: {
            Authorization: localStorage.getItem('jwt'),
          },
        }
      )
      .then((response) => {
        for (const i in response.data) {
          if (response.data[i].amount >= 0) {
            setBalancesPositive1((prev) => [
              ...prev,
              {
                amount: response.data[i].amount,
                category: response.data[i].category.name,
              },
            ]);
          } else {
            setBalancesNegative1((prev) => [
              ...prev,
              {
                amount: Math.abs(response.data[i].amount),
                category: response.data[i].category.name,
              },
            ]);
          }
        }

        setLoading(false);
        setStatsReady(true);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  const handleStat2Changed = () => {
    setBalancesNegative2({ nodes: [], links: [] });
    setBalancesPositive2({ nodes: [], links: [] });
    axios
      .post(
        '/api/stats',
        {
          selectedAccounts: [accountSelected],
          selectedCategories: categories.map((category) => category.id),
          initDate: new Date('2020-01-01').getTime(),
          endDate: new Date('2020-12-31').getTime(),
        },
        {
          headers: {
            Authorization: localStorage.getItem('jwt'),
          },
        }
      )
      .then((response) => {
        let balancesPositive: Balances2 = {
          nodes: [
            {
              name: accounts.find((account) => account.id === accountSelected)!.name,
            },
          ],
          links: [],
        };

        let balancesNegative: Balances2 = {
          nodes: [
            {
              name: accounts.find((account) => account.id === accountSelected)!.name,
            },
          ],
          links: [],
        };

        for (const i in response.data) {
          if (response.data[i].amount > 0) {
            const nodes = balancesPositive.nodes;
            const links = balancesPositive.links;
            if (nodes.some((node) => node.name === response.data[i].category.name)) {
              links.push({
                source: 0,
                target: nodes.indexOf(nodes.find((node) => node.name === response.data[i].category.name)!),
                value: response.data[i].amount,
              });
            } else {
              nodes.push({
                name: response.data[i].category.name,
              });
              links.push({
                source: 0,
                target: nodes.length - 1,
                value: response.data[i].amount,
              });
            }

            balancesPositive = {
              nodes: nodes,
              links: links,
            };
          } else if (response.data[i].amount < 0) {
            const nodes = balancesNegative.nodes;
            const links = balancesNegative.links;
            if (nodes.some((node) => node.name === response.data[i].category.name)) {
              links.push({
                source: 0,
                target: nodes.indexOf(nodes.find((node) => node.name === response.data[i].category.name)!),
                value: Math.abs(response.data[i].amount),
              });
            } else {
              nodes.push({
                name: response.data[i].category.name,
              });
              links.push({
                source: 0,
                target: nodes.length - 1,
                value: Math.abs(response.data[i].amount),
              });
            }

            balancesNegative = {
              nodes: nodes,
              links: links,
            };
          }
        }

        setBalancesPositive2(balancesPositive);
        setBalancesNegative2(balancesNegative);

        setLoading(false);
        setStatsReady(true);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  const handleStat3Changed = () => {
    setBalances3([]);
    axios
      .get('/api/accounts/full', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        const transactions = [];
        const initDate = new Date('2020-01-01').getTime();
        const endDate = new Date('2020-12-31').getTime();
        const interval = Math.floor((endDate - initDate) / 4);
        const dateValues = [initDate, initDate + interval, initDate + interval * 2, initDate + interval * 3, endDate];

        for (let i in response.data) {
          for (let j in response.data[i].transactions) {
            if (
              response.data[i].transactions[j].category.id === categoriesPositionSelected[0] ||
              response.data[i].transactions[j].category.id === categoriesPositionSelected[1]
            ) {
              transactions.push(response.data[i].transactions[j]);
            }
          }
        }

        let data = dateValues.map((value) => {
          const dataObj: any = {
            date: value,
          };
          dataObj[categories.find((category) => category.id === categoriesPositionSelected[0])!.name] = 0;
          dataObj[categories.find((category) => category.id === categoriesPositionSelected[1])!.name] = 0;
          return dataObj;
        });

        for (let i = 1; i < data.length; i++) {
          data[i][categories.find((category) => category.id === categoriesPositionSelected[0])!.name] =
            data[i - 1][categories.find((category) => category.id === categoriesPositionSelected[0])!.name];
          data[i][categories.find((category) => category.id === categoriesPositionSelected[1])!.name] =
            data[i - 1][categories.find((category) => category.id === categoriesPositionSelected[1])!.name];
          for (let j in transactions) {
            if (transactions[j].date >= data[i - 1].date && transactions[j].date < data[i].date) {
              data[i][transactions[j].category.name] += transactions[j].amount;
            }
          }
        }

        data = data.map((item) => ({
          ...item,
          date: new Date(item.date)
            .toLocaleDateString('es-ES', {
              year: '2-digit',
              month: '2-digit',
              day: '2-digit',
            })
            .replace(/\//g, '-'),
        }));

        setBalances3(data);
        setCategoriesComparison([
          categories.find((category) => category.id === categoriesPositionSelected[0])!.name,
          categories.find((category) => category.id === categoriesPositionSelected[1])!.name,
        ]);

        setLoading(false);
        setStatsReady(true);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  const handleSelectCategory = (category: string, index: number) => {
    setCategoriesPositionSelected((prevState) => {
      return { ...prevState, [index]: category };
    });
  };

  useEffect(() => {
    setStatsReady(false);
    setCategoriesPositionSelected({
      0: '',
      1: '',
    });
    handleStatTypeChanged();
  }, [statType]);

  useEffect(() => {
    setStatsReady(false);
    if (statType === '1') {
      if (accountsSelected.length > 0 && categoriesSelected.length > 0) {
        setLoading(true);
        handleStat1Changed();
      }
    } else if (statType === '2') {
      if (accountSelected !== '') {
        setLoading(true);
        handleStat2Changed();
      }
    } else if (statType === '3') {
      console.log(categoriesPositionSelected);
      if (
        categoriesPositionSelected[0] !== '' &&
        categoriesPositionSelected[1] !== '' &&
        categoriesPositionSelected[0] !== categoriesPositionSelected[1]
      ) {
        setLoading(true);
        handleStat3Changed();
      }
    }
  }, [accountSelected, accountsSelected, categoriesSelected, categoriesPositionSelected]);

  return (
    <div className={styles.statsPage}>
      <div className={styles.titleWrapper}>
        <span className={styles.title}>Stats</span>
      </div>
      <div className={styles.body}>
        <div className={styles.typeStats}>
          <div className={styles.selectWrapper}>
            <InputSelect placeholder='Select a type of stat' options={statTypes} onChange={setStatType} />
          </div>
        </div>
        <div className={styles.configWrapper}>
          {statType === '1' && (
            <>
              <div className={styles.selectWrapper}>
                <InputSelect
                  placeholder='Select at least an account'
                  options={accounts}
                  onChange={() => {}}
                  multiple={true}
                  onChangeMultiple={setAccountsSelected}
                />
              </div>
              <div className={styles.selectWrapper}>
                <InputSelect
                  placeholder='Select at least a category'
                  options={categories}
                  onChange={() => {}}
                  multiple={true}
                  onChangeMultiple={setCategoriesSelected}
                />
              </div>
            </>
          )}
          {statType === '2' && (
            <div className={styles.selectWrapper}>
              <InputSelect placeholder='Select an account' options={accounts} onChange={setAccountSelected} />
            </div>
          )}
          {statType === '3' && (
            <>
              <div className={styles.selectWrapper}>
                <InputSelect
                  placeholder='Select a category'
                  options={categories}
                  onChange={(value) => handleSelectCategory(value, 0)}
                />
              </div>
              <div className={styles.selectWrapper}>
                <InputSelect
                  placeholder='Select a category'
                  options={categories}
                  onChange={(value) => handleSelectCategory(value, 1)}
                />
              </div>
            </>
          )}
        </div>
        <div className={styles.contentWrapper}>
          {loading && <Loading />}

          {statType === '1' && statsReady && (
            <CategoriesBalances balancesPositive={balancesPositive1} balancesNegative={balancesNegative1} />
          )}
          {statType === '2' && statsReady && (
            <AccountBalances balancesPositive={balancesPositive2} balancesNegative={balancesNegative2} />
          )}
          {statType === '3' && statsReady && <CategoriesComparison balances={balances3} categories={categoriesComparison} />}
        </div>
      </div>
    </div>
  );
};

export default StatsPage;
