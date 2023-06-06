import Image from 'next/image';
import styles from './InputText.module.css';

type Props = {
  iconPath: string;
  onChange: (value: string) => void;
  type?: string;
};

const InputText = ({ iconPath, onChange, type = 'text' }: Props) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value);
  };

  return (
    <div className={styles.inputText}>
      <Image className={styles.icon} src={iconPath} alt='Login image' height={30} width={30} />
      <input className={styles.input} type={type} onChange={handleChange} />
    </div>
  );
};

export default InputText;
