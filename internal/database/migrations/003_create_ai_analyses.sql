-- Migration: 003_create_ai_analyses
-- Description: Cria tabela para armazenar análises de IA dos manuscritos
-- Author: Sprint 3-4
-- Date: 2025-10-31

-- Create ai_analyses table
CREATE TABLE IF NOT EXISTS ai_analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manuscript_id UUID NOT NULL,
    
    -- Classificação de gênero
    genre VARCHAR(100) NOT NULL,
    sub_genre VARCHAR(100),
    genre_confidence DECIMAL(5,4),
    
    -- Análise de tom
    tone VARCHAR(100) NOT NULL,
    tone_score DECIMAL(5,4),
    
    -- Elementos especiais
    equation_percentage DECIMAL(5,4) DEFAULT 0.0,
    code_percentage DECIMAL(5,4) DEFAULT 0.0,
    table_count INTEGER DEFAULT 0,
    image_count INTEGER DEFAULT 0,
    
    -- Recomendações de pipeline
    recommended_pipeline VARCHAR(50) NOT NULL,
    pipeline_confidence DECIMAL(5,4),
    pipeline_reason TEXT,
    
    -- Recomendações de fontes
    recommended_body_font VARCHAR(100),
    recommended_title_font VARCHAR(100),
    recommended_mono_font VARCHAR(100),
    font_rationale TEXT,
    
    -- Metadados
    analyzed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    tokens_used INTEGER DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Foreign Keys
    CONSTRAINT fk_ai_analyses_manuscript
        FOREIGN KEY (manuscript_id)
        REFERENCES manuscripts(id)
        ON DELETE CASCADE
);

-- Indexes para performance
CREATE INDEX idx_ai_analyses_manuscript_id ON ai_analyses(manuscript_id);
CREATE INDEX idx_ai_analyses_genre ON ai_analyses(genre);
CREATE INDEX idx_ai_analyses_analyzed_at ON ai_analyses(analyzed_at DESC);

-- Trigger para atualizar updated_at
CREATE OR REPLACE FUNCTION update_ai_analyses_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_ai_analyses_updated_at
    BEFORE UPDATE ON ai_analyses
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_analyses_updated_at();

-- Comentários
COMMENT ON TABLE ai_analyses IS 'Armazena análises de IA de manuscritos (gênero, tom, elementos especiais, recomendações)';
COMMENT ON COLUMN ai_analyses.genre_confidence IS 'Confiança da classificação de gênero (0.0 a 1.0)';
COMMENT ON COLUMN ai_analyses.equation_percentage IS 'Percentual de equações no conteúdo';
COMMENT ON COLUMN ai_analyses.code_percentage IS 'Percentual de código no conteúdo';
COMMENT ON COLUMN ai_analyses.tokens_used IS 'Tokens consumidos da API OpenAI nesta análise';
