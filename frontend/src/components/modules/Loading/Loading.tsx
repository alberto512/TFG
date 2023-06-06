import Image from 'next/image';
import DotsLoading from '@elements/DotsLoading/DotsLoading';
import styles from './Loading.module.css';

const Loading = () => {
  return (
    <div className={styles.loading}>
      <DotsLoading />
      <Image src={'/svg/loading.svg'} alt='Loading' height={500} width={500} />
    </div>
  );
};

export default Loading;
