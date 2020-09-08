package lib

import (
	"errors"
)

// ----------------------------------------------------------------
func (this *Mlrmap) Has(key *string) bool {
	return this.findEntry(key) != nil
}

func (this *Mlrmap) findEntry(key *string) *mlrmapEntry {
	if this.keysToEntries != nil {
		return this.keysToEntries[*key]
	} else {
		for pe := this.Head; pe != nil; pe = pe.Next {
			if *pe.Key == *key {
				return pe
			}
		}
		return nil
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) Put(key *string, value *Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value)
		if this.Head == nil {
			this.Head = pe
			this.Tail = pe
		} else {
			pe.Prev = this.Tail
			pe.Next = nil
			this.Tail.Next = pe
			this.Tail = pe
		}
		if this.keysToEntries != nil {
			this.keysToEntries[*key] = pe
		}
		this.FieldCount++
	} else {
		copy := *value
		pe.Value = &copy
	}
}

// ----------------------------------------------------------------
// E.g. '$name[1]["foo"] = "bar"'
// The key is "name" and the indices are [1, "foo"].
// See also indexed-lvalues.md.
func (this *Mlrmap) PutIndexed(key *string, indices []*Mlrval, rvalue *Mlrval) error {
	mlrval := this.Get(key)
	if mlrval == nil {
		mapval := MlrvalEmptyMap()
		this.Put(key, &mapval)
		return mapval.PutIndexed(indices, rvalue)
	} else if mlrval.IsAbsent() {
		return errors.New("Value [\"" + *key + "\"] is not a maaaaaap.")
	} else if !mlrval.IsMap() {
		return errors.New("Value [\"" + *key + "\"] is not a map.")
	} else {
		return mlrval.PutIndexed(indices, rvalue)
	}
}

// E.g. '$*["foo"][1] = "bar"'
// The indices are ["foo", 1].
// See also indexed-lvalues.md.
//
// This is a Mlrmap (from string to Mlrval) so we handle the first level of
// indexing here, then pass the remaining indices to the Mlrval at the desired
// slot.
func (this *Mlrmap) PutIndexedKeyless(indices []*Mlrval, rvalue *Mlrval) error {
	n := len(indices)
	if n == 0 { // mlr put '$* = {"a":1, "b":2}'
		if !rvalue.IsMap() {
			return errors.New("Cannot assign non-map to existing map; got " + rvalue.GetTypeName() + ".")
		}
		*this = *rvalue.mapval // xxx needs deepcopy
		return nil
	}

	baseIndex := indices[0]
	if !baseIndex.IsString() {
		skey := baseIndex.String()
		return errors.New("Non-string key " + skey) // xxx needs better wording
	}
	baseKey := baseIndex.printrep

	if n == 1 {
		this.Put(&baseKey, rvalue) // E.g. mlr put '$*["a"] = 3'
		return nil
	}

	baseValue := this.Get(&baseKey)
	if baseValue == nil {
		baseValue := MlrvalEmptyMap()
		this.Put(&baseKey, &baseValue)
		return baseValue.PutIndexed(indices[1:], rvalue)
	} else if !baseValue.IsMap() {
		return errors.New("Value [\"" + baseKey + "\"] is not a map; got " + baseValue.GetTypeName() + ".")
	} else {
		return baseValue.PutIndexed(indices[1:], rvalue)
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) Prepend(key *string, value *Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value)
		if this.Tail == nil {
			this.Head = pe
			this.Tail = pe
		} else {
			pe.Prev = nil
			pe.Next = this.Head
			this.Head.Prev = pe
			this.Head = pe
		}
		if this.keysToEntries != nil {
			this.keysToEntries[*key] = pe
		}
		this.FieldCount++
	} else {
		copy := *value
		pe.Value = &copy
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) Get(key *string) *Mlrval {
	pe := this.findEntry(key)
	if pe == nil {
		return nil
	} else {
		return pe.Value
	}
	return nil
}

func (this *Mlrmap) Clear() {
	this.FieldCount = 0
	// Assuming everything unreferenced is getting GC'ed by the Go runtime
	this.Head = nil
	this.Tail = nil
}

// ----------------------------------------------------------------
// TODO: needs to be a deepcopy -- Mlrval needs its own Copy method.
func (this *Mlrmap) Copy() *Mlrmap {
	that := NewMlrmapMaybeHashed(this.isHashed())
	for pe := this.Head; pe != nil; pe = pe.Next {
		that.Put(pe.Key, pe.Value)
	}
	return that
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (this *Mlrmap) Remove(key *string) bool {
	pe := this.findEntry(key)
	if pe == nil {
		return false
	} else {
		this.unlink(pe)
		return true
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) unlink(pe *mlrmapEntry) {
	if pe == this.Head {
		if pe == this.Tail {
			this.Head = nil
			this.Tail = nil
		} else {
			this.Head = pe.Next
			pe.Next.Prev = nil
		}
	} else {
		pe.Prev.Next = pe.Next
		if pe == this.Tail {
			this.Tail = pe.Prev
		} else {
			pe.Next.Prev = pe.Prev
		}
	}
	if this.keysToEntries != nil {
		delete(this.keysToEntries, *pe.Key)
	}
	this.FieldCount--
}

//mlrmapEntry* lrec_put_after(Mlrmap* prec, mlrmapEntry* pd, char* key, char* value, char free_flags) {
//	mlrmapEntry* pe = lrec_find_entry(prec, key);
//
//	if (pe != NULL) { // Overwrite
//		if (pe->free_flags & FREE_ENTRY_VALUE) {
//			free(pe->value);
//		}
//		pe->value = value;
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//		if (free_flags & FREE_ENTRY_VALUE)
//			pe->free_flags |= FREE_ENTRY_VALUE;
//	} else { // Insert after specified entry
//		pe = mlr_malloc_or_die(sizeof(mlrmapEntry));
//		pe->key         = key;
//		pe->value       = value;
//		pe->free_flags  = free_flags;
//		pe->quote_flags = 0;
//
//		if (pd->Next == NULL) { // Append at end of list
//			pd->Next = pe;
//			pe->Prev = pd;
//			pe->Next = NULL;
//			prec->Tail = pe;
//
//		} else {
//			mlrmapEntry* pf = pd->Next;
//			pd->Next = pe;
//			pf->Prev = pe;
//			pe->Prev = pd;
//			pe->Next = pf;
//		}
//
//		prec->field_count++;
//	}
//	return pe;
//}

//char* lrec_get_ext(Mlrmap* prec, char* key, mlrmapEntry** ppentry) {
//	mlrmapEntry* pe = lrec_find_entry(prec, key);
//	if (pe != NULL) {
//		*ppentry = pe;
//		return pe->value;
//	} else {
//		*ppentry = NULL;;
//		return NULL;
//	}
//}

//// ----------------------------------------------------------------
//mlrmapEntry* lrec_get_pair_by_position(Mlrmap* prec, int position) { // 1-up not 0-up
//	if (position <= 0 || position > prec->field_count) {
//		return NULL;
//	}
//	int sought_index = position - 1;
//	int found_index = 0;
//	mlrmapEntry* pe = NULL;
//	for (
//		found_index = 0, pe = prec->Head;
//		pe != NULL;
//		found_index++, pe = pe->Next
//	) {
//		if (found_index == sought_index) {
//			return pe;
//		}
//	}
//	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//		MLR_GLOBALS.bargv0, __FILE__, __LINE__);
//	exit(1);
//}

//char* lrec_get_key_by_position(Mlrmap* prec, int position) { // 1-up not 0-up
//	mlrmapEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		return NULL;
//	} else {
//		return pe->key;
//	}
//}

//char* lrec_get_value_by_position(Mlrmap* prec, int position) { // 1-up not 0-up
//	mlrmapEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		return NULL;
//	} else {
//		return pe->value;
//	}
//}

//// ----------------------------------------------------------------
//void lrec_remove_by_position(Mlrmap* prec, int position) { // 1-up not 0-up
//	mlrmapEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//
//	free(pe);
//}

// Before:
//   "x" => "3"
//   "y" => "4"  <-- pold
//   "z" => "5"  <-- pnew
//
// Rename y to z
//
// After:
//   "x" => "3"
//   "z" => "4"
//
//void lrec_rename(Mlrmap* prec, char* old_key, char* new_key, int new_needs_freeing) {
//
//	mlrmapEntry* pold = lrec_find_entry(prec, old_key);
//	if (pold != NULL) {
//		mlrmapEntry* pnew = lrec_find_entry(prec, new_key);
//
//		if (pnew == NULL) { // E.g. rename "x" to "y" when "y" is not present
//			if (pold->free_flags & FREE_ENTRY_KEY) {
//				free(pold->key);
//				pold->key = new_key;
//				if (!new_needs_freeing)
//					pold->free_flags &= ~FREE_ENTRY_KEY;
//			} else {
//				pold->key = new_key;
//				if (new_needs_freeing)
//					pold->free_flags |=  FREE_ENTRY_KEY;
//			}
//
//		} else { // E.g. rename "x" to "y" when "y" is already present
//			if (pnew->free_flags & FREE_ENTRY_VALUE) {
//				free(pnew->value);
//			}
//			if (pold->free_flags & FREE_ENTRY_KEY) {
//				free(pold->key);
//				pold->free_flags &= ~FREE_ENTRY_KEY;
//			}
//			pold->key = new_key;
//			if (new_needs_freeing)
//				pold->free_flags |=  FREE_ENTRY_KEY;
//			else
//				pold->free_flags &= ~FREE_ENTRY_KEY;
//			lrec_unlink(prec, pnew);
//			free(pnew);
//		}
//	}
//}

// Cases:
// 1. Rename field at position 3 from "x" to "y when "y" does not exist elsewhere in the srec
// 2. Rename field at position 3 from "x" to "y when "y" does     exist elsewhere in the srec
// Note: position is 1-up not 0-up
//void  lrec_rename_at_position(Mlrmap* prec, int position, char* new_key, int new_needs_freeing){
//	mlrmapEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		if (new_needs_freeing) {
//			free(new_key);
//		}
//		return;
//	}
//
//	mlrmapEntry* pother = lrec_find_entry(prec, new_key);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	pe->key = new_key;
//	if (new_needs_freeing) {
//		pe->free_flags |= FREE_ENTRY_KEY;
//	} else {
//		pe->free_flags &= ~FREE_ENTRY_KEY;
//	}
//	if (pother != NULL) {
//		lrec_unlink(prec, pother);
//		free(pother);
//	}
//}

//// ----------------------------------------------------------------
//void lrec_move_to_head(Mlrmap* prec, char* key) {
//	mlrmapEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//	lrec_link_at_head(prec, pe);
//}

//void lrec_move_to_tail(Mlrmap* prec, char* key) {
//	mlrmapEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//	lrec_link_at_tail(prec, pe);
//}

// ----------------------------------------------------------------
// Simply rename the first (at most) n positions where n is the length of pnames.
//
// Possible complications:
//
// * pnames itself contains duplicates -- we require this as invariant-check
//   from the caller since (for performance) we don't want to check this on every
//   record processed.
//
// * pnames has length less than the current record and one of the new names
//   becomes a clash with an existing name.
//
//   Example:
//   - Input record has names "a,b,c,d,e".
//   - pnames is "d,x,f"
//   - We then construct the invalid "d,x,f,d,e" -- we need to detect and unset
//     the second 'd' field.

//void  lrec_label(Mlrmap* prec, slls_t* pnames_as_list, hss_t* pnames_as_set) {
//	mlrmapEntry* pe = prec->Head;
//	sllse_t* pn = pnames_as_list->Head;
//
//	// Process the labels list
//	for ( ; pe != NULL && pn != NULL; pe = pe->Next, pn = pn->Next) {
//		char* new_name = pn->value;
//
//		if (pe->free_flags & FREE_ENTRY_KEY) {
//			free(pe->key);
//		}
//		pe->key = mlr_strdup_or_die(new_name);;
//		pe->free_flags |= FREE_ENTRY_KEY;
//	}
//
//	// Process the remaining fields in the record beyond those affected by the new-labels list
//	for ( ; pe != NULL; ) {
//		char* name = pe->key;
//		if (hss_has(pnames_as_set, name)) {
//			mlrmapEntry* Next = pe->Next;
//			if (pe->free_flags & FREE_ENTRY_KEY) {
//				free(pe->key);
//			}
//			if (pe->free_flags & FREE_ENTRY_VALUE) {
//				free(pe->value);
//			}
//			lrec_unlink(prec, pe);
//			free(pe);
//			pe = Next;
//		} else {
//			pe = pe->Next;
//		}
//	}
//}

//// ----------------------------------------------------------------
//void lrece_update_value(mlrmapEntry* pe, char* new_value, int new_needs_freeing) {
//	if (pe == NULL) {
//		return;
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//	pe->value = new_value;
//	if (new_needs_freeing)
//		pe->free_flags |= FREE_ENTRY_VALUE;
//	else
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//}

//// ----------------------------------------------------------------
//static void lrec_link_at_head(Mlrmap* prec, mlrmapEntry* pe) {
//
//	if (prec->Head == NULL) {
//		pe->Prev   = NULL;
//		pe->Next   = NULL;
//		prec->Head = pe;
//		prec->Tail = pe;
//	} else {
//		// [b,c,d] + a
//		pe->Prev   = NULL;
//		pe->Next   = prec->Head;
//		prec->Head->Prev = pe;
//		prec->Head = pe;
//	}
//	prec->field_count++;
//}

//static void lrec_link_at_tail(Mlrmap* prec, mlrmapEntry* pe) {
//
//	if (prec->Head == NULL) {
//		pe->Prev   = NULL;
//		pe->Next   = NULL;
//		prec->Head = pe;
//		prec->Tail = pe;
//	} else {
//		pe->Prev   = prec->Tail;
//		pe->Next   = NULL;
//		prec->Tail->Next = pe;
//		prec->Tail = pe;
//	}
//	prec->field_count++;
//}

//// ----------------------------------------------------------------
//Mlrmap* lrec_literal_1(char* k1, char* v1) {
//	Mlrmap* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	return prec;
//}

//Mlrmap* lrec_literal_2(char* k1, char* v1, char* k2, char* v2) {
//	Mlrmap* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	return prec;
//}

//Mlrmap* lrec_literal_3(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3) {
//	Mlrmap* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	lrec_put(prec, k3, v3, NO_FREE);
//	return prec;
//}

//Mlrmap* lrec_literal_4(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3, char* k4, char* v4) {
//	Mlrmap* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	lrec_put(prec, k3, v3, NO_FREE);
//	lrec_put(prec, k4, v4, NO_FREE);
//	return prec;
//}
