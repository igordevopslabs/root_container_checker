# Root Container Checker
Essa é uma simples aplicação para identificar containers executando com privilégios root dentro de um cluster K8S.
Atualmente, essa aplicação pode ser executada somente dentro de um cluster K8S, isso por que ela se autentica no ambeinte através de uma serviceAccount, não sendo possivel passar passar um kubeconfig.

## Como instalar?

Você pode instalar ela no seu cluster a partir do [helmfile](https://github.com/helmfile/helmfile)
```bash
helmfile -f deploy/helmfile.yaml -e demo apply
```

Você também pode gerar o helm template da aplicação com o comando ```helmfile -f deploy/helmfile.yaml -e demo tempalte``` e aplica-la diretamente via ```kubectl apply```.

O helmfile foi escolhido para simplificar o gerenciamento dos charts, bem como a manutenção e atualização de valores.

### Período
Para evitar uma sobrecarga na api do K8S, essa aplicação possui um "ticker" para ser executada a cada 8hrs, ou seja, a cada ciclo de 8hrs, a aplicação bate na api do K8S para realizar a consulta de quantos containers estão executando com usuário root.

### Autenticação e Autorização
Essa aplicação possui autenticação e autirização restrita, de maneira que ela só consegue interagir coma api do K8S através de uma Service Account com escopo bem limitado:

```yaml
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```


## Exportação de métricas
Essa aplicação foi pensada para exportar métricas prometheus que contabilizam a quantidade de containers rodando como root no ambiente, trazendo assim um mapeamento a respeito dos containers que podem se tornar risco em caso de comprometimento. 

![image](https://user-images.githubusercontent.com/73206099/236629778-0a98dd87-d707-4140-b38a-95d47b021105.png)

## Dashboard
A partir das métricas é possivel elaborar um dashboard para realizar o mapeamento do desses containers:
![image](https://user-images.githubusercontent.com/73206099/236630319-67ff5738-b482-46ec-a7f0-5dfeb709562d.png)
