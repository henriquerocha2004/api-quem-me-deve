--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13
-- Dumped by pg_dump version 17.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cancel_info; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cancel_info (
    id character(26) NOT NULL,
    reason text NOT NULL,
    cancel_date timestamp without time zone,
    cancelled_by character(26) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    debt_id character(26) NOT NULL
);


--
-- Name: debts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.debts (
    id character(26) NOT NULL,
    description text NOT NULL,
    total_value numeric(12,2) NOT NULL,
    due_date timestamp without time zone,
    installments_quantity integer NOT NULL,
    debt_date timestamp without time zone,
    status character varying(255) NOT NULL,
    user_client_id character(26) NOT NULL,
    product_ids character(26)[] DEFAULT '{}'::bpchar[],
    service_ids character(26)[] DEFAULT '{}'::bpchar[],
    finished_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: installments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.installments (
    id character(26) NOT NULL,
    description text NOT NULL,
    value numeric(12,2) NOT NULL,
    due_date timestamp without time zone,
    deb_date timestamp without time zone,
    status character varying(255) NOT NULL,
    payment_date timestamp without time zone,
    payment_method character varying(255),
    number integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    debt_id character(26) NOT NULL
);


--
-- Name: reversal_info; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reversal_info (
    id character(26) NOT NULL,
    reason text NOT NULL,
    reversal_date timestamp without time zone,
    reversed_by character(26) NOT NULL,
    reversed_installment_qtd integer DEFAULT 0 NOT NULL,
    cancelled_installment_qtd integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    debt_id character(26) NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: cancel_info cancel_info_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cancel_info
    ADD CONSTRAINT cancel_info_pkey PRIMARY KEY (id);


--
-- Name: debts debts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.debts
    ADD CONSTRAINT debts_pkey PRIMARY KEY (id);


--
-- Name: installments installments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installments
    ADD CONSTRAINT installments_pkey PRIMARY KEY (id);


--
-- Name: reversal_info reversal_info_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reversal_info
    ADD CONSTRAINT reversal_info_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: idx_cancel_info_cancel_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cancel_info_cancel_date ON public.cancel_info USING btree (cancel_date);


--
-- Name: idx_cancel_info_cancelled_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cancel_info_cancelled_by ON public.cancel_info USING btree (cancelled_by);


--
-- Name: idx_debts_debt_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_debts_debt_date ON public.debts USING btree (debt_date);


--
-- Name: idx_debts_finished_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_debts_finished_at ON public.debts USING btree (finished_at);


--
-- Name: idx_debts_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_debts_status ON public.debts USING btree (status);


--
-- Name: idx_debts_user_client_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_debts_user_client_id ON public.debts USING btree (user_client_id);


--
-- Name: idx_installments_deb_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_installments_deb_date ON public.installments USING btree (deb_date);


--
-- Name: idx_installments_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_installments_number ON public.installments USING btree (number);


--
-- Name: idx_installments_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_installments_status ON public.installments USING btree (status);


--
-- Name: idx_reversal_info_reversal_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reversal_info_reversal_date ON public.reversal_info USING btree (reversal_date);


--
-- Name: idx_reversal_info_reversed_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reversal_info_reversed_by ON public.reversal_info USING btree (reversed_by);


--
-- Name: cancel_info cancel_info_debt_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cancel_info
    ADD CONSTRAINT cancel_info_debt_id_fkey FOREIGN KEY (debt_id) REFERENCES public.debts(id);


--
-- Name: installments installments_debt_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installments
    ADD CONSTRAINT installments_debt_id_fkey FOREIGN KEY (debt_id) REFERENCES public.debts(id);


--
-- Name: reversal_info reversal_info_debt_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reversal_info
    ADD CONSTRAINT reversal_info_debt_id_fkey FOREIGN KEY (debt_id) REFERENCES public.debts(id);


--
-- PostgreSQL database dump complete
--

