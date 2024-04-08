package entities

type KeyAttribute struct {
	UserID                            string `xorm:"not null pk VARCHAR(36) user_id"`
	KekSalt                           string `xorm:"not null TEXT kek_salt"`
	EncryptedKey                      string `xorm:"not null TEXT encrypted_key"`
	KeyDecryptionNonce                string `xorm:"not null TEXT key_decryption_nonce"`
	PublicKey                         string `xorm:"not null TEXT public_key"`
	EncryptedSecretKey                string `xorm:"not null TEXT encrypted_secret_key"`
	SecretKeyDecryptionNonce          string `xomr:"not null TEXT secret_key_decryption_none"`
	MasterKeyEncryptedWithRecoveryKey string `xorm:"not null TEXT master_key_encrypted_with_recovery_key"`
	MasterKeyDecryptionNonce          string `xorm:"not null TEXT master_key_decryption_nonce"`
	RecoveryKeyEncryptedWithMasterKey string `xorm:"not null TEXT recovery_key_encrypted_with_master_key"`
	RecoveryKeyDecryptionNonce        string `xorm:"not null TEXT recovery_key_decryption_nonce"`
	MemLimit                          int    `xorm:"not null TEXT mem_limit"`
	OpsLimit                          int    `xorm:"not null TEXT ops_limit"`
}

func (*KeyAttribute) TableName() string {
	return "key_attributes"
}
