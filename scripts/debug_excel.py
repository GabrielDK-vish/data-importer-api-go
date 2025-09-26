#!/usr/bin/env python3
"""
Script para debug do arquivo Excel
"""
import pandas as pd
import sys

def debug_excel(file_path):
    try:
        # Ler arquivo Excel
        df = pd.read_excel(file_path, sheet_name=0)
        
        print(f"📊 Arquivo: {file_path}")
        print(f"📏 Dimensões: {df.shape[0]} linhas x {df.shape[1]} colunas")
        print(f"📋 Colunas encontradas:")
        
        for i, col in enumerate(df.columns):
            print(f"  {i:2d}: {col}")
        
        # Verificar colunas obrigatórias
        required = ['PartnerId', 'CustomerId', 'ProductId', 'UsageDate', 'Quantity', 'UnitPrice']
        print(f"\n🔍 Verificação de colunas obrigatórias:")
        
        missing = []
        for col in required:
            if col in df.columns:
                print(f"  ✅ {col}")
            else:
                print(f"  ❌ {col} - FALTANDO")
                missing.append(col)
        
        if missing:
            print(f"\n❌ Colunas obrigatórias não encontradas: {missing}")
            print("💡 Sugestões de mapeamento:")
            for col in df.columns:
                if 'partner' in col.lower() or 'id' in col.lower():
                    print(f"  '{col}' -> PartnerId")
                elif 'customer' in col.lower():
                    print(f"  '{col}' -> CustomerId")
                elif 'product' in col.lower():
                    print(f"  '{col}' -> ProductId")
                elif 'usage' in col.lower() or 'date' in col.lower():
                    print(f"  '{col}' -> UsageDate")
                elif 'quantity' in col.lower():
                    print(f"  '{col}' -> Quantity")
                elif 'price' in col.lower() or 'unit' in col.lower():
                    print(f"  '{col}' -> UnitPrice")
        else:
            print(f"\n✅ Todas as colunas obrigatórias foram encontradas!")
        
        # Mostrar primeiras linhas
        print(f"\n📄 Primeiras 3 linhas:")
        print(df.head(3).to_string())
        
    except Exception as e:
        print(f"❌ Erro ao processar arquivo: {e}")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        debug_excel(sys.argv[1])
    else:
        debug_excel("../Reconfile fornecedores.xlsx")
